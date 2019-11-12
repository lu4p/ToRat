package client

import (
	"os/exec"
	"runtime"
)

func runCmd(cmd string, powershell bool) []byte {
	var osshell string
	if cmd == "" {
		return []byte("Error: No command to execute!")
	}
	var osshellargs []string
	if runtime.GOOS == "linux" {
		osshell = "/bin/sh"
		osshellargs = []string{"-c", cmd}

	} else if runtime.GOOS == "windows" {
		if powershell {
			osshell = "powershell"
			osshellargs = []string{"-Command", cmd}

		} else {
			osshell = "cmd"
			osshellargs = []string{"/C", cmd}
		}
	} else if runtime.GOOS == "darwin" {
		// TODO: Add right strings for Mac OSX
		osshell = ""
		osshellargs = []string{"", cmd}
	}
	execcmd := exec.Command(osshell, osshellargs...)
	cmdout, err := execcmd.Output()
	if err != nil {
		return []byte("err")
	} else if cmdout == nil {
		return []byte("no output!")
	}
	return cmdout

}
