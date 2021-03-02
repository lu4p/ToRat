package client

import (
	"crypto/rsa"
	"crypto/x509"
	_ "embed" // used for embedding the cert
	"encoding/pem"
	"log"
)

const (
	serverPort string = ":1337"
)

type server struct {
	cert   *x509.Certificate
	pubKey *rsa.PublicKey
	addr   string
	domain string
}

var s server

//go:embed cert.pem
var serverCert []byte

// initServer returns a struct with the cert, domain, pubkey, and address
// for dialing the source server's tor address
func initServer() {
	serverBlock, _ := pem.Decode(serverCert)

	cert, err := x509.ParseCertificate(serverBlock.Bytes)
	if err != nil {
		panic(err)
	}

	domain := cert.DNSNames[0]
	log.Println("[initServer] Initialized server cert")

	s = server{
		cert:   cert,
		addr:   domain + serverPort,
		pubKey: cert.PublicKey.(*rsa.PublicKey),
		domain: domain,
	}
}
