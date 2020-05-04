package config

import (
	"strconv"
)

// CreateVolumeConfig - values needed to create a volume
type CreateVolumeConfig struct {
	GIGS int
}

// CreateVolume - Get config for create volume
func CreateVolume() *CreateVolumeConfig {
	volume := GetEnv("VOLUME_SIZE")
	if volume != nil {
		volume, _ := strconv.Atoi(*volume)
		return &CreateVolumeConfig{volume}
	}

	return nil
}
