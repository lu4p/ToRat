package client

import (
	"github.com/pierrre/archivefile/zip"
)

// GetVer gets the major version of the current installed OSX
func GetVer() (int, error) {

}

// CheckElevate checks whether the current process has administrator
// privileges
func CheckElevate() bool {

}

// TODO: Remove Zipdir because .zip files are insecure and can be exploited
// Zipdir archives files to zip and sends them to server
func (c *connection) Zipdir(path string) error {
	progress := func(archivePath string) {
	}
	return zip.Archive(path, c.Conn, progress)
}
