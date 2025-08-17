package k8s

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/curve25519"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"

	"vpnaas-backend/internal/config"
	"vpnaas-backend/internal/metrics"
	"vpnaas-backend/internal/models"
)

// VPNManager handles VPN pod lifecycle and configuration
type VPNManager struct {
	clientset *kubernetes.Clientset
	namespace string
}

// WireGuardKeys represents a pair of WireGuard keys
type WireGuardKeys struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

// NewVPNManager creates a new VPN manager
func NewVPNManager(clientset *kubernetes.Clientset) *VPNManager {
	namespace := config.GetString("k8s.namespace")
	if namespace == "" {
		namespace = "vpnaas"
	}

	return &VPNManager{
		clientset: clientset,
		namespace: namespace,
	}
}

// CreateUserVPN creates a VPN pod for a user
func (vm *VPNManager) CreateUserVPN(ctx context.Context, user *models.User) error {
	// Generate WireGuard keys
	keys, err := vm.generateWireGuardKeys()
	if err != nil {
		return fmt.Errorf("failed to generate WireGuard keys: %v", err)
	}

	user.PublicKey = keys.PublicKey
	user.PrivateKey = keys.PrivateKey

	// Generate WireGuard configuration
	config, err := vm.generateWireGuardConfig(user, keys)
	if err != nil {
		return fmt.Errorf("failed to generate WireGuard config: %v", err)
	}

	user.ConfigData = config

	// Create Kubernetes pod
	pod, err := vm.createVPNPod(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create VPN pod: %v", err)
	}

	user.PodName = pod.Name
	user.PodIP = pod.Status.PodIP

	logrus.Infof("Created VPN pod %s for user %s", pod.Name, user.Username)

	// Update metrics
	metrics.IncrementConnections()

	return nil
}

// DeleteUserVPN deletes a VPN pod for a user
func (vm *VPNManager) DeleteUserVPN(ctx context.Context, user *models.User) error {
	if user.PodName == "" {
		return nil
	}

	err := vm.clientset.CoreV1().Pods(vm.namespace).Delete(ctx, user.PodName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete VPN pod: %v", err)
	}

	logrus.Infof("Deleted VPN pod %s for user %s", user.PodName, user.Username)

	return nil
}

// GetPodStatus returns the status of a VPN pod
func (vm *VPNManager) GetPodStatus(ctx context.Context, podName string) (string, error) {
	pod, err := vm.clientset.CoreV1().Pods(vm.namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(pod.Status.Phase), nil
}

// UpdatePodMetrics updates pod-related metrics
func (vm *VPNManager) UpdatePodMetrics(ctx context.Context) error {
	pods, err := vm.clientset.CoreV1().Pods(vm.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=vpnaas,component=vpn",
	})
	if err != nil {
		return err
	}

	running, failed, pending := 0, 0, 0
	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			running++
		case corev1.PodFailed:
			failed++
		case corev1.PodPending:
			pending++
		}
	}

	metrics.UpdatePodMetrics(running, failed, pending)
	return nil
}

// generateWireGuardKeys generates a new WireGuard key pair
func (vm *VPNManager) generateWireGuardKeys() (*WireGuardKeys, error) {
	privateKey := make([]byte, 32)
	if _, err := rand.Read(privateKey); err != nil {
		return nil, err
	}

	var publicKey [32]byte
	curve25519.ScalarBaseMult(&publicKey, (*[32]byte)(privateKey))

	return &WireGuardKeys{
		PrivateKey: base64.StdEncoding.EncodeToString(privateKey),
		PublicKey:  base64.StdEncoding.EncodeToString(publicKey[:]),
	}, nil
}

// generateWireGuardConfig generates WireGuard configuration for a user
func (vm *VPNManager) generateWireGuardConfig(user *models.User, keys *WireGuardKeys) (string, error) {
	// Generate a unique IP for the user (10.0.0.x)
	userIP := fmt.Sprintf("10.0.0.%d", len(user.ID)%254+1)

	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s/24
ListenPort = %s
PostUp = iptables -A FORWARD -i %%i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i %%i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE

[Peer]
PublicKey = %s
AllowedIPs = 0.0.0.0/0
Endpoint = %s:%s
PersistentKeepalive = 25
`,
		keys.PrivateKey,
		userIP,
		config.GetString("vpn.wireguard_port"),
		keys.PublicKey,
		config.GetString("vpn.endpoint"),
		config.GetString("vpn.wireguard_port"),
	)

	return config, nil
}

// createVPNPod creates a Kubernetes pod for VPN
func (vm *VPNManager) createVPNPod(ctx context.Context, user *models.User) (*corev1.Pod, error) {
	podName := fmt.Sprintf("vpn-%s", user.ID)

	// Create ConfigMap for WireGuard configuration
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("vpn-config-%s", user.ID),
			Namespace: vm.namespace,
			Labels: map[string]string{
				"app":       "vpnaas",
				"component": "vpn",
				"user":      user.ID,
			},
		},
		Data: map[string]string{
			"wg0.conf": user.ConfigData,
		},
	}

	_, err := vm.clientset.CoreV1().ConfigMaps(vm.namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create ConfigMap: %v", err)
	}

	// Create the pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: vm.namespace,
			Labels: map[string]string{
				"app":       "vpnaas",
				"component": "vpn",
				"user":      user.ID,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "wireguard",
					Image: config.GetString("vpn.image"),
					Ports: []corev1.ContainerPort{
						{
							Name:          "wireguard",
							ContainerPort: int32(config.GetInt("vpn.wireguard_port")),
							Protocol:      corev1.ProtocolUDP,
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config",
							MountPath: "/config/wg_confs",
						},
					},
					SecurityContext: &corev1.SecurityContext{
						Capabilities: &corev1.Capabilities{
							Add: []corev1.Capability{
								"NET_ADMIN",
								"SYS_MODULE",
							},
						},
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse(config.GetString("vpn.pod_cpu_request")),
							corev1.ResourceMemory: resource.MustParse(config.GetString("vpn.pod_memory_request")),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse(config.GetString("vpn.pod_cpu_limit")),
							corev1.ResourceMemory: resource.MustParse(config.GetString("vpn.pod_memory_limit")),
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: configMap.Name,
							},
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	createdPod, err := vm.clientset.CoreV1().Pods(vm.namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	// Wait for pod to be ready
	err = vm.waitForPodReady(ctx, createdPod.Name)
	if err != nil {
		return nil, err
	}

	return createdPod, nil
}

// waitForPodReady waits for a pod to be ready
func (vm *VPNManager) waitForPodReady(ctx context.Context, podName string) error {
	return retry.OnError(retry.DefaultRetry, func(err error) bool {
		return true
	}, func() error {
		pod, err := vm.clientset.CoreV1().Pods(vm.namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if pod.Status.Phase == corev1.PodRunning {
			return nil
		}

		return fmt.Errorf("pod not ready yet")
	})
}
