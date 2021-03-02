package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/cretz/bine/tor"
)

func main() {
	err := mainErr()
	if err != nil {
		log.Fatalln(err)
	}
}

func mainErr() error {
	_, hsPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}
	t, err := tor.Start(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("unable to start Tor: %v", err)
	}
	defer t.Close()

	onion, err := t.Listen(context.Background(), &tor.ListenConf{
		RemotePorts: []int{80},
		Key:         hsPriv,
	})
	if err != nil {
		return err
	}

	log.Println(onion.ID + ".onion")

	defer onion.Close()

	if err := ioutil.WriteFile("hs_private", hsPriv, 0o666); err != nil {
		return err
	}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	template := x509.Certificate{
		IsCA:                  true,
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Acme Co"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{onion.ID + ".onion"},
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &rsaKey.PublicKey, rsaKey)
	if err != nil {
		return err
	}

	if err := writePemKey("../torat_client/cert.pem", "CERTIFICATE", cert); err != nil {
		return err
	}

	return writePemKey(
		"priv_key.pem",
		"RSA PRIVATE KEY",
		x509.MarshalPKCS1PrivateKey(rsaKey),
	)
}

func writePemKey(name, typ string, bytes []byte) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	defer f.Close()

	privBlock := pem.Block{
		Type:  typ,
		Bytes: bytes,
	}

	return pem.Encode(f, &privBlock)
}
