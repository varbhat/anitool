package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
)

func AES256Encrypt(hexkey string, hexiv string, text string) string {
	key, _ := hex.DecodeString(hexkey)
	iv, _ := hex.DecodeString(hexiv)
	return aes256([]byte(text), key, iv, aes.BlockSize)
}

func aes256(plaintext []byte, key []byte, iv []byte, blockSize int) string {
	bPlaintext := pKCS5Padding(plaintext, blockSize, len(plaintext))
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	encryptedString := base64.StdEncoding.EncodeToString(ciphertext)
	return encryptedString
}

func pKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
