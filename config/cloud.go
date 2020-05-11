package config

type CloudProvider string

const (
	CloudProviderDigitalOcean CloudProvider = "digitalocean"
)

// CloudConfig - values needed to run anything against a cloud provider
type CloudConfig struct {
	Provider CloudProvider
	APIKey   string
}

// GetCloudConfig - Get config for perform operation on cloud
func GetCloudConfig() *CloudConfig {
	volume := GetEnv("CLOUD_DIGITALOCEAN_API_KEY", "cloud.digitalocean.api_key")
	if volume != nil {
		return &CloudConfig{CloudProviderDigitalOcean, *volume}
	}

	return nil
}
