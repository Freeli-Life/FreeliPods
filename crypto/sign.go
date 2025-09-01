package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"time"
)

func LoadPrivateKey(path string) (crypto.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
	}

	return key, nil
}

func SignRegistrationData(privateKey crypto.PrivateKey, domain, username string, salt, sigKey, encKey []byte) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	var message []byte
	message = append(message, []byte(domain)...)
	message = append(message, []byte(username)...)
	message = append(message, salt...)
	message = append(message, []byte(timestamp)...)
	message = append(message, sigKey...)
	message = append(message, encKey...)

	hashed := sha256.Sum256(message)

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an RSA key, signing is not supported")
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	return signature, nil
}