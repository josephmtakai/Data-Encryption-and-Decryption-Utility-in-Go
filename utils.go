package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

// SavePrivateKey saves the RSA private key to a file
func SavePrivateKey(filename string, key *rsa.PrivateKey) error {
	privKeyFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer privKeyFile.Close()

	pem.Encode(privKeyFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return nil
}

// LoadPrivateKey loads an RSA private key from a file
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	privKeyData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privKeyData)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

// Additional RSA functions can be added here for encryption/decryption
