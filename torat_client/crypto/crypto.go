package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"io"
	"log"
	"os"

	"github.com/lu4p/ToRat/shared"
)

// GenRandString generate a random string
func GenRandString() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln("Couldn't generate hostname:", err)
	}

	return string(base64.RawURLEncoding.EncodeToString(b))
}

// genHostname generates the Hostname of the machine
func genHostname(PubKey *rsa.PublicKey) shared.EncAsym {
	hostname := GenRandString()
	return encAsym([]byte(hostname), PubKey)
}

// GetHostname returns the encrypted Hostname
// if Hostname is not set a new Hostname is generated
func GetHostname(path string, PubKey *rsa.PublicKey) shared.EncAsym {
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

func getEncodedFile(path string) (shared.EncAsym, error) {
	f, err := os.Open(path)
	if err != nil {
		return shared.EncAsym{}, err
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	var encAsym shared.EncAsym
	decoder.Decode(&encAsym)
	return encAsym, nil
}

// encAsym encrypts data using RSA + AES to the publickey
// of the server
func encAsym(data []byte, pubKey *rsa.PublicKey) shared.EncAsym {
	aeskey := genAesKey()
	encKey := encRsa(aeskey, pubKey)
	encData := encAes(data, aeskey)
	return shared.EncAsym{
		EncAesKey: encKey,
		EncData:   encData,
	}
}

// encodeToFile encodes data using gob and writes the result to path
func encodeToFile(data shared.EncAsym, path string) error {
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
