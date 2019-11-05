package client

import (
	"os"
	"path/filepath"
)

var (
	Path         = filepath.Join(os.ExpandEnv("$HOME"), ".cache", "MacSafe")
	PathExe      = filepath.Join(Path, "MacSafe")
	HostnamePath = filepath.Join(Path, "token")
)
