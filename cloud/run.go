package cloud

import (
	"github.com/levibostian/purslane/config"
)

// CreateVolume Creates a volume from the cloud provider the user wants to run.
func CreateVolume(cloudConfig *config.CloudConfig, config *config.CreateVolumeConfig) *CreatedVolume {
	cloudProvider := getCloudProvider(cloudConfig)

	return cloudProvider.createVolume(config)
}

// CreateServer - create server
func CreateServer(coreConfig *config.CoreConfig, cloudConfig *config.CloudConfig, createServerConfig *config.CreateServerConfig, createdVolume *CreatedVolume) *CreatedServer {
	cloudProvider := getCloudProvider(cloudConfig)

	return cloudProvider.createServer(coreConfig, createServerConfig, createdVolume)
}

func DeleteVolume(cloudConfig *config.CloudConfig, createdVolume *CreatedVolume) {
	getCloudProvider(cloudConfig).deleteVolume(createdVolume)
}

func DeleteServer(cloudConfig *config.CloudConfig, createdServer *CreatedServer) {
	getCloudProvider(cloudConfig).deleteServer(createdServer)
}

func getCloudProvider(cloudConfig *config.CloudConfig) Cloud {
	// Right now, we are only allowing DigitalOcean as a provider so, let's just create it.

	//if cloudConfig.Provider == config.CloudProviderDigitalOcean {
	return GetDigitalOceanCloud(cloudConfig)
	//}
}
