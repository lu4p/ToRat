package shared

import (
	"fmt"
	"os/exec"
	"runtime"
)

func RunCmd(cmd string, powershell bool) ([]byte, error) {
	if cmd == "" {
		return nil, fmt.Errorf("no command to execute (empty command)")
	}

	var (
		osshell     string
		osshellargs []string
	)

	switch runtime.GOOS {
	case "linux":
		osshell = "/bin/sh"
		osshellargs = []string{"-c", cmd}
	case "windows":
		if powershell {
			osshell = "powershell"
			osshellargs = []string{"-Command", cmd}

		} else {
			osshell = "cmd"
			osshellargs = []string{"/C", cmd}
		}
	case "darwin":
		osshell = "/bin/sh"
		osshellargs = []string{"-c", cmd}
	}

	execcmd := exec.Command(osshell, osshellargs...)
	cmdout, err := execcmd.Output()
	if err != nil {
		return nil, err
	}

	return cmdout, nil
}
