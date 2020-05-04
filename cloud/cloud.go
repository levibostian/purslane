package cloud

import (
	"github.com/levibostian/purslane/config"
)

// Cloud provider
type Cloud interface {
	createVolume(config.CreateVolumeConfig)
}
