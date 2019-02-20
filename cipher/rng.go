package cipher

import (
	"crypto/rand"
)

func GenerateSalt(len int) []byte {
	salt := make([]byte, len)
	rand.Read(salt)
	return salt
}

func GenerateIV() []byte {
	iv := make([]byte, 16)
	rand.Read(iv)
	return iv
}
