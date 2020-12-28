package cloud

import (
	"github.com/levibostian/purslane/config"
)

func GetCloudProvider(coreConfig *config.CoreConfig, cloudConfig *config.CloudConfig) CloudProvider {
	// Right now, we are only allowing DigitalOcean as a provider so, let's just create it.

	//if cloudConfig.Provider == config.CloudProviderDigitalOcean {
	return GetDigitalOceanCloud(coreConfig, cloudConfig)
	//}
}
