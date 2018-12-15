package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func Hash(value string) string {
	h := sha256.Sum256([]byte(value))
	b := []byte(h[:])
	return base64.StdEncoding.EncodeToString(b)
}

func padKey(key []byte) []byte {
	k := len(key)
	if k < 16 {
		return append(key, bytes.Repeat([]byte{byte(k)}, 16-k)...)
	} else if k < 24 {
		return append(key, bytes.Repeat([]byte{byte(k)}, 24-k)...)
	} else if k < 32 {
		return append(key, bytes.Repeat([]byte{byte(k)}, 32-k)...)
	}

	return nil
}

func Encrypt(key []byte, data []byte) (cipherText []byte, err error) {
	key = padKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText = make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)

	return
}

func Decrypt(key []byte, data []byte) (dec []byte, err error) {
	key = padKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(data) < aes.BlockSize {
		err = errors.New("Block size is too short!")
		return
	}

	iv := data[:aes.BlockSize]
	dec = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(dec, dec)

	return
}
