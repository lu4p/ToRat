package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/lu4p/binclude"
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

//go:generate binclude

func initServer() {
	certName := binclude.Include("../keygen/cert.pem")
	serverCert, err := BinFS.ReadFile(certName)
	if err != nil {
		panic(err)
	}

	serverBlock, _ := pem.Decode(serverCert)

	cert, err := x509.ParseCertificate(serverBlock.Bytes)
	if err != nil {
		panic(err)
	}

	domain := cert.DNSNames[0]
	log.Println("Domain:", domain)

	s = server{
		cert:   cert,
		addr:   domain + serverPort,
		pubKey: cert.PublicKey.(*rsa.PublicKey),
		domain: domain,
	}
}
