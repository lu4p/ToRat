package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/lu4p/ToRat/shared"
)

var privateKey = loadPrivateKey()

func loadPrivateKey() *rsa.PrivateKey {
	key, err := ioutil.ReadFile("../../keygen/priv_key.pem")
	if err != nil {
		log.Fatalln("Cannot read PrivateKey:", err)
	}

	block, _ := pem.Decode(key)
	if block == nil {
		log.Fatalln("PrivateKey is not PEM format:", err)
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalln("Could not parse PrivateKey:", err)
	}

	return priv
}

// DecRsa decrypts RSA-encrypted data
func DecRsa(encData []byte) ([]byte, error) {
	rng := rand.Reader
	decData, err := rsa.DecryptOAEP(sha256.New(), rng, privateKey, encData, nil)
	if err != nil {
		log.Println("[!] Rsa:", err)
		return nil, err
	}

	return decData, nil
}

// DecAes decrypts data encrypted with AES
func DecAes(encData []byte, aeskey []byte) ([]byte, error) {
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, encData[:12], encData[12:], nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// DecAsym decypts asymetric encryption (4096 bit RSA + AES)
func DecAsym(encData shared.EncAsym) ([]byte, error) {
	aeskey, err := DecRsa(encData.EncAesKey)
	if err != nil {
		return nil, err
	}

	return DecAes(encData.EncData, aeskey)
}
