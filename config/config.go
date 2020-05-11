package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/levibostian/purslane/ui"
	"github.com/levibostian/purslane/util"
)

type CoreConfig struct {
	PublicSSHKeyFingerprint string
	PublicSSHKey            string
}

func GetCoreConfig() *CoreConfig {
	publicSSHKeyPath := GetEnv("PUBLIC_SSH_KEY_PATH", "public_ssh_key_path")
	if publicSSHKeyPath != nil {
		info, err := os.Stat(*publicSSHKeyPath)
		if os.IsNotExist(err) {
			ui.Abort(fmt.Sprintf("Public SSH key given with path, %s, does not exist", *publicSSHKeyPath))
		}
		if info.IsDir() {
			ui.Abort(fmt.Sprintf("Public SSH key given with path, %s, is a directory and not a public SSH key file.", *publicSSHKeyPath))
		}
		content, err := ioutil.ReadFile(*publicSSHKeyPath)
		ui.HandleError(err)
		publicSSHKey := string(content)

		fingerprint := util.GetSSHFootprint(publicSSHKey)

		return &CoreConfig{fingerprint, publicSSHKey}
	}

	return nil
}
