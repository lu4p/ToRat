package client

import (
	"os"
	"path/filepath"
)

var (
	Path         = filepath.Join(os.ExpandEnv("$HOME"), ".cache", "libssh")
	PathExe      = filepath.Join(Path, "libssh")
	HostnamePath = filepath.Join(Path, "token")
)
