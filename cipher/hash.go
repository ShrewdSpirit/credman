package cipher

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"hash"

	"golang.org/x/crypto/scrypt"
)

type ScryptDifficulty int

const (
	ScryptDifficultyEasy ScryptDifficulty = 8
	ScryptDifficultyNorm ScryptDifficulty = 16
	ScryptDifficultyHard ScryptDifficulty = 24
)

func ScryptWithSalt(input string, diff ScryptDifficulty, keylen int) (hash, salt []byte) {
	salt = GenerateSalt(32)
	hash = Scrypt(input, diff, keylen, salt)
	return
}

func Scrypt(input string, r ScryptDifficulty, keylen int, salt []byte) (hash []byte) {
	hash, _ = scrypt.Key([]byte(input), salt, 65536, int(r), 1, keylen)
	return
}

func GenerateSalt(len int) (salt []byte) {
	salt = make([]byte, len)
	rand.Read(salt)
	return
}

func GenerateIV() (iv []byte) {
	iv = make([]byte, aes.BlockSize)
	rand.Read(iv)
	return
}

// BuildKeys builds required keys for encryption and decryption.
// If inputSalt is empty, it will generate one and fill salt output.
func BuildKeys(key string, diff ScryptDifficulty, inputSalt []byte) (aesKey, hmacKey, check, salt []byte) {
	var masterHash []byte

	if inputSalt == nil {
		masterHash, salt = ScryptWithSalt(key, diff, 128)
	} else {
		masterHash = Scrypt(key, diff, 128, inputSalt)
	}

	aesKey = masterHash[:32]
	hmacKey = masterHash[32:64]
	check = masterHash[64:] // this is double the sizes

	return
}

func BuildHmac(dataLen uint64, encryptedData, hmacKey, iv []byte) (h hash.Hash) {
	h = hmac.New(sha256.New, hmacKey)
	h.Write(iv)
	binary.Write(h, binary.LittleEndian, dataLen)
	h.Write(encryptedData)

	return
}
