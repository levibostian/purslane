package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/levibostian/Purslane/ui"
)

// Go by default does not handle ~/ in a file path. So, we try to take a string and expand it to the full file path.
func GetFullFilePath(filePath string) string {
	// Use strings.HasPrefix so we don't match paths like
	// "/something/~/something/"
	if strings.HasPrefix(filePath, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir

		filePath = filepath.Join(dir, filePath[2:])
	}

	ui.Debug("Util get full file path. Result: %s", filePath)

	return filePath
}

func GetFileContents(path string, fileDescribe string) []byte {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		ui.Abort(fmt.Sprintf("%s file at path, %s, does not exist", fileDescribe, path))
	}
	if info.IsDir() {
		ui.Abort(fmt.Sprintf("%s file at path, %s, is a directory and not a file.", fileDescribe, path))
	}
	content, err := ioutil.ReadFile(path)
	ui.HandleError(err)

	return content
}
