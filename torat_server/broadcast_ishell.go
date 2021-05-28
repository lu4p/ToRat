package server

import (
	"github.com/abiosoft/ishell"
	"path/filepath"
)

// BroadcastShell Broadcast shell
func (bc broadcast) BroadcastShell() {

	serverFileCompleter := func([]string) []string {
		files, err := filepath.Glob("*")
		if err != nil {
			return nil
		}
		return files
	}

	// Set shell and get working dir
	shell := ishell.New()

	shell.SetPrompt(yellow("[Broadcast]$ "))

	commands := []*ishell.Cmd{
		{
			Name: "shred",
			Func: bc.Shred,
			Help: "remove a path by overwriting it with random data then removing it: usage shred <path>",
		},
		{
			Name:      "up",
			Func:      bc.Upload,
			Help:      "upload a file from the cwd of the Server to cwd of the client: usage up <file>",
			Completer: serverFileCompleter,
		},
		{
			Name: "escape",
			Func: bc.runCommand,
			Help: "escape a command and run it natively on client",
		},
		{
			Name: "exit",
			Func: func(c *ishell.Context) {
				c.Stop()
				shell.Close()
			},
			Help: "exit broadcast shell",
		},
	}

	for _, c := range commands {
		shell.AddCmd(c)
	}

	shell.NotFound(bc.runCommand)
	shell.Run()
}

func (bc *broadcast) Shred(c *ishell.Context) {
	for _, ac := range activeClients {
		ac.Shred(c)
	}
}

func (bc *broadcast) Upload(c *ishell.Context) {
	for _, ac := range activeClients {
		ac.Upload(c)
	}
}

func (bc *broadcast) runCommand(c *ishell.Context) {
	for _, ac := range activeClients {
		ac.runCommand(c)
	}
}
