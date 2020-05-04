package cloud

import (
	"github.com/levibostian/purslane/config"
)

// CreateVolume Creates a volume from the cloud provider the user wants to run.
func CreateVolume(config config.CreateVolumeConfig) {
	digitaloceanCloud.createVolume(config)
}
