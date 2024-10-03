package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"

	"encoding/base64"

	"errors"
	"io/ioutil"
	"strings"
)

// GenerateRSAKeyPair generates a new RSA key pair
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

// EncryptMessage encrypts a plain text message using AES
func EncryptMessage(plainText string) (string, error) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(plainText))

	// Encode key and ciphertext for storage
	encodedKey := base64.StdEncoding.EncodeToString(key)
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedKey + ":" + encodedCiphertext, nil
}

// DecryptMessage decrypts an encrypted message using AES
func DecryptMessage(encryptedData string) (string, error) {
	parts := strings.Split(encryptedData, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid encrypted data format")
	}
	key, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// EncryptFile encrypts a file using AES
func EncryptFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return EncryptMessage(string(data))
}

// DecryptFile decrypts a file using AES
func DecryptFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return DecryptMessage(string(data))
}
