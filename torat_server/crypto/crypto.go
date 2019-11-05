package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

var privateKey = loadPrivateKey()

func loadPrivateKey() *rsa.PrivateKey {
	key, err := ioutil.ReadFile("key.pem")
	if err != nil {
		log.Println("err read", err)
		return nil
	}
	block, _ := pem.Decode(key)
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Could not parse rsakey", err)
		return nil
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

// DecAsym decypts asymetric encryption (2048 bit RSA + AES)
func DecAsym(encData []byte) ([]byte, error) {
	if len(encData) < 256 {
		return nil, errors.New("Unsufficent AesKey length")
	}
	encAeskey := encData[:256]
	encContent := encData[256:]
	log.Println("before rsa")
	aeskey, err := DecRsa(encAeskey)
	if err != nil {
		return nil, err
	}
	log.Println("after rsa")
	return DecAes(encContent, aeskey)
}
