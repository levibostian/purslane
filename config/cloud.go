package config

import (
	"strconv"
)

// CloudConfig - values needed to run anything against a cloud provider
type CloudConfig struct {
	API_KEY string
}

// Create - Get config for perform operation on cloud
func Cloud() *CloudConfig {
	volume := GetEnv("VOLUME_SIZE")
	if volume != nil {
		volume, _ := strconv.Atoi(*volume)
		return &CreateVolumeConfig{volume}
	}

	return nil
}
