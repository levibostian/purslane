package util

import (
	"github.com/levibostian/Purslane/ui"
	"golang.org/x/crypto/ssh"
)

func GetSSHFootprint(publicKey string) string {
	pubKeyBytes := []byte(publicKey)

	pk, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	ui.HandleError(err)

	fingerprint := ssh.FingerprintLegacyMD5(pk)

	return fingerprint
}
