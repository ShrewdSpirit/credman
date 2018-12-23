package cipher

import (
	"golang.org/x/crypto/scrypt"
)

type ScryptDifficulty int

const (
	ScryptDifficultyEasy ScryptDifficulty = 0
	ScryptDifficultyNorm ScryptDifficulty = 1
	ScryptDifficultyHard ScryptDifficulty = 2
)

func HashScrypt(input string, diff ScryptDifficulty, hashlen int) (hash, salt []byte) {
	salt = GenerateSalt(32)
	hash = HashScryptSalt(input, diff, hashlen, salt)
	return
}

func HashScryptSalt(input string, diff ScryptDifficulty, hashlen int, salt []byte) (hash []byte) {
	r := 8
	switch diff {
	case ScryptDifficultyEasy:
		r = 8
	case ScryptDifficultyNorm:
		r = 16
	case ScryptDifficultyHard:
		r = 24
	}

	hash, _ = scrypt.Key([]byte(input), salt, 65536, r, 1, hashlen)
	return
}
