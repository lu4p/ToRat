package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
	mathrand "math/rand"
)

// CertToPubKey convert a X509 PEM encoded certicate to RSA PublicKey
func CertToPubKey(certPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(certPEM))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return rsaPublicKey, nil
}

// GenRandString generate a random string
func GenRandString() string {
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	length := 16
	byte := make([]byte, length)
	for i := 0; i < length; i++ {
		num := mathrand.Intn(len(all))
		byte[i] = all[num]
	}
	return string(byte)
}

// SetHostname Sets the Hostname of the machine to a
// random string with the length of 16, encrypts the outcome and
// writes it to Disk
func SetHostname(path string, PubKey *rsa.PublicKey) error {
	hostname := GenRandString()
	return EnctoFile([]byte(hostname), path, PubKey)
}

// GetHostname returns the encrypted Hostname
// if Hostname is not set a new Hostname is generated
func GetHostname(path string, PubKey *rsa.PublicKey) []byte {
	encHostname, err := ioutil.ReadFile(path)
	if err != nil {
		if SetHostname(path, PubKey) == nil {
			encHostname, err = ioutil.ReadFile(path)
			if err != nil {
				return nil
			}
			return encHostname
		}
		return nil
	}
	return encHostname
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

// EnctoFile encrypts data using RSA + AES to the publickey
// of the server and writes the encrypted data to disk
func EnctoFile(data []byte, path string, PubKey *rsa.PublicKey) error {
	aeskey, err := genAesKey()
	if err != nil {
		return err
	}
	encKey := encRsa(aeskey, PubKey)
	encData, err := encAes(data, aeskey)
	if err != nil {
		return err
	}
	enc := append(encKey, encData...)
	err = ioutil.WriteFile(path, enc, 0666)
	if err != nil {
		return err
	}
	return nil

}

// genAesKey generates a 256bit AES Key
func genAesKey() ([]byte, error) {
	AesKey := make([]byte, 32)
	_, err := rand.Read(AesKey)
	if err != nil {
		return nil, err
	}
	return AesKey, nil
}

func encAes(data []byte, AesKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(AesKey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	encData := append(nonce, ciphertext...)
	return encData, nil
}
