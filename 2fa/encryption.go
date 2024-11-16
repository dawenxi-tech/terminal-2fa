package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

var (
	encryptKey = `AgomCgr2qRIm5hc0`
	_parsedKey []byte
)

func init() {
	key, _ := base64.StdEncoding.DecodeString(encryptKey)
	if len(key) == 0 {
		key = []byte(encryptKey)
	}
	for len(key) < 16 {
		key = append(key, key...)
	}
	_parsedKey = key[:16]
}

func encrypt(msg string) (string, error) {
	return encryptMessage(_parsedKey, msg)
}

func decrypt(msg string) (string, error) {
	return decryptMessage(_parsedKey, msg)
}

// aes encryption
// https://gist.github.com/fracasula/38aa1a4e7481f9cedfa78a0cdd5f1865
func encryptMessage(key []byte, message string) (string, error) {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}
	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decryptMessage(key []byte, message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}
	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}
