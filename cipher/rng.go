package cipher

import (
	"math/rand"
	"time"
)

var randSource = rand.NewSource(time.Now().UnixNano())
var rng = rand.New(randSource)

func RandRange(min, max int64) int64 {
	return min + rng.Int63()%max
}

func GenerateSalt(len int) []byte {
	salt := make([]byte, len)
	rng.Read(salt)
	return salt
}

func GenerateIV() []byte {
	iv := make([]byte, 16)
	rng.Read(iv)
	return iv
}
