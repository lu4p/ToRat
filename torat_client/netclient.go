// +build tor

package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/cretz/bine/process/embedded"
	"github.com/cretz/bine/tor"
)

func connect(dialer *tor.Dialer) (net.Conn, error) {
	log.Println("connecting to", s.addr)
	conn, err := dialer.Dial("tcp", s.addr)
	if err != nil {
		return nil, err
	}

	log.Println("connect")
	caPool := x509.NewCertPool()
	caPool.AddCert(s.cert)

	config := tls.Config{RootCAs: caPool, ServerName: s.domain}
	tlsconn := tls.Client(conn, &config)

	return tlsconn, nil
}

// NetClient start tor and invoke connect
func NetClient() {
	log.Println("NetClient")
	initServer()
	var conf tor.StartConf
	conf = tor.StartConf{ProcessCreator: embedded.NewCreator()}

	t, err := tor.Start(nil, &conf)
	if err != nil {
		log.Println("[!] Tor could not be started:", err)
		return
	}

	api := new(API)
	rpc_err := rpc.Register(api)
	if rpc_err != nil {
		log.Fatal(rpc_err)
	}

	defer t.Close()
	dialer, _ := t.Dialer(nil, nil)

	for {
		conn, err := connect(dialer)
		if err != nil {
			log.Println("Could not connect:", err)
			time.Sleep(10 * time.Second)
			continue
		}
		rpc.ServeConn(conn)
	}
}
