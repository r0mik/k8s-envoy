package config

import (
	"github.com/spf13/viper"
)

// Load initializes configuration from environment variables and config files
func Load() error {
	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("debug", false)
	viper.SetDefault("vpn.wireguard_port", "51820")
	viper.SetDefault("vpn.pod_cpu_limit", "100m")
	viper.SetDefault("vpn.pod_memory_limit", "128Mi")
	viper.SetDefault("vpn.pod_cpu_request", "50m")
	viper.SetDefault("vpn.pod_memory_request", "64Mi")
	viper.SetDefault("vpn.image", "linuxserver/wireguard:latest")
	viper.SetDefault("k8s.namespace", "vpnaas")
	viper.SetDefault("k8s.pod_labels", map[string]string{
		"app": "vpnaas",
		"component": "vpn",
	})

	// Read from environment variables
	viper.SetEnvPrefix("VPNAAS")
	viper.AutomaticEnv()

	// Read from config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/vpnaas")

	// Ignore config file not found error
	viper.ReadInConfig()

	return nil
}

// GetString returns a string configuration value
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns an integer configuration value
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a boolean configuration value
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetStringMap returns a string map configuration value
func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}
