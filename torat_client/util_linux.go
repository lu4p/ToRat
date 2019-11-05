package client

import (
	"errors"

	"github.com/pierrre/archivefile/zip"
)

// GetVer gets the major version of the current installed linux
func GetVer() (int, error) {
	return 0, errors.New("Could not get Ver")
}

// CheckElevate checks whether the current process has administrator
// privileges
func CheckElevate() bool {
	return false
}

// TODO: Remove Zipdir because .zip files are insecure and can be exploited
// Zipdir archives files to zip and sends them to server
func (c *connection) Zipdir(path string) error {
	progress := func(archivePath string) {
	}
	return zip.Archive(path, c.Conn, progress)
}
