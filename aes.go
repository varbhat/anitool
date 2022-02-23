package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func AES256Encrypt(hexkey string, hexiv string, text string) string {
	key, err := hex.DecodeString(hexkey)
	if err != nil {
		panic(err)
	}
	iv, err := hex.DecodeString(hexiv)
	if err != nil {
		panic(err)
	}
	return aes256encrypt([]byte(text), key, iv, aes.BlockSize)
}

func aes256encrypt(plaintext []byte, key []byte, iv []byte, blockSize int) string {
	bPlaintext := pKCS5Padding(plaintext, blockSize, len(plaintext))
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	encryptedString := base64.StdEncoding.EncodeToString(ciphertext)
	return encryptedString
}

func aes256decrypt(plaintext string, key []byte, iv []byte, blocksize int) string {
	bPlaintext, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		fmt.Println(err)
	}
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	ciphertext = pkCS7unpad(ciphertext, blocksize)
	return string(ciphertext)
}

func pKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkCS7unpad(padded []byte, size int) []byte {
	if len(padded)%size != 0 {
		return nil
	}

	bufLen := len(padded) - int(padded[len(padded)-1])
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf
}
