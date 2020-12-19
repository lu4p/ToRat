package client

import (
	"errors"
)

// GetVer gets the major version of the current installed linux
func GetVer() (int, error) {
	return 0, errors.New("could not get Version")
}

// CheckElevate checks whether the current process has administrator
// privileges
func CheckElevate() bool {
	return false
}
