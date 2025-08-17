package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// User metrics
	TotalUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vpnaas_total_users",
		Help: "Total number of VPN users",
	})

	ActiveUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vpnaas_active_users",
		Help: "Number of active VPN users",
	})

	InactiveUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vpnaas_inactive_users",
		Help: "Number of inactive VPN users",
	})

	SuspendedUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vpnaas_suspended_users",
		Help: "Number of suspended VPN users",
	})

	// VPN pod metrics
	VPNPodsRunning = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vpnaas_vpn_pods_running",
		Help: "Number of running VPN pods",
	})

	VPNPodsFailed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vpnaas_vpn_pods_failed",
		Help: "Number of failed VPN pods",
	})

	VPNPodsPending = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vpnaas_vpn_pods_pending",
		Help: "Number of pending VPN pods",
	})

	// Connection metrics
	TotalConnections = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vpnaas_total_connections",
		Help: "Total number of VPN connections",
	})

	ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vpnaas_active_connections",
		Help: "Number of currently active VPN connections",
	})

	// Data usage metrics
	TotalDataUsage = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vpnaas_total_data_usage_bytes",
		Help: "Total data usage in bytes",
	})

	DataUsagePerUser = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vpnaas_user_data_usage_bytes",
		Help: "Data usage per user in bytes",
	}, []string{"user_id", "username"})

	// API metrics
	APIRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vpnaas_api_requests_total",
		Help: "Total number of API requests",
	}, []string{"method", "endpoint", "status"})

	APIRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "vpnaas_api_request_duration_seconds",
		Help:    "API request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint"})

	// Error metrics
	ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vpnaas_errors_total",
		Help: "Total number of errors",
	}, []string{"type", "component"})
)

// Init initializes the metrics
func Init() {
	// Initialize all metrics to 0
	TotalUsers.Set(0)
	ActiveUsers.Set(0)
	InactiveUsers.Set(0)
	SuspendedUsers.Set(0)
	VPNPodsRunning.Set(0)
	VPNPodsFailed.Set(0)
	VPNPodsPending.Set(0)
	ActiveConnections.Set(0)
	TotalDataUsage.Add(0)
}

// UpdateUserMetrics updates user-related metrics
func UpdateUserMetrics(total, active, inactive, suspended int) {
	TotalUsers.Set(float64(total))
	ActiveUsers.Set(float64(active))
	InactiveUsers.Set(float64(inactive))
	SuspendedUsers.Set(float64(suspended))
}

// UpdatePodMetrics updates pod-related metrics
func UpdatePodMetrics(running, failed, pending int) {
	VPNPodsRunning.Set(float64(running))
	VPNPodsFailed.Set(float64(failed))
	VPNPodsPending.Set(float64(pending))
}

// IncrementConnections increments the connection counter
func IncrementConnections() {
	TotalConnections.Inc()
}

// SetActiveConnections sets the number of active connections
func SetActiveConnections(count int) {
	ActiveConnections.Set(float64(count))
}

// AddDataUsage adds data usage to the total counter
func AddDataUsage(bytes int64) {
	TotalDataUsage.Add(float64(bytes))
}

// SetUserDataUsage sets data usage for a specific user
func SetUserDataUsage(userID, username string, bytes int64) {
	DataUsagePerUser.WithLabelValues(userID, username).Set(float64(bytes))
}

// RecordAPIRequest records an API request
func RecordAPIRequest(method, endpoint, status string) {
	APIRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}

// RecordAPIRequestDuration records API request duration
func RecordAPIRequestDuration(method, endpoint string, duration float64) {
	APIRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// RecordError records an error
func RecordError(errorType, component string) {
	ErrorsTotal.WithLabelValues(errorType, component).Inc()
}
