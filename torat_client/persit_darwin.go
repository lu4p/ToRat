package client

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	Path         = filepath.Join(os.ExpandEnv("$HOME"), ".cache", "MacSafe")
	PathExe      = filepath.Join(Path, "MacSafe")
	HostnamePath = filepath.Join(Path, "token")
)

// GetVer gets the major version of the current installed OSX
func GetVer() (int, error) {
	// TODO: Implement
	return 0, errors.New("not implemented")
}

// CheckElevate checks whether the current process has administrator
// privileges
func CheckElevate() bool {
	// TODO: implement
	return false
}

// Persist makes sure that the executable is run after a reboot
func Persist(path string) {
	elevated := CheckElevate()
	if elevated {
		persistAdmin(path)
		return
	}

	persistUser(path)
}

// persistAdmin persistence using admin privileges
func persistAdmin(path string) {
}

// persistUser persistence using user privileges
func persistUser(path string) {
}
