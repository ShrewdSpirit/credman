package cipher

import (
	"crypto/sha256"
	"errors"
)

const bufSize = 16 * 1024
const hmacSize = sha256.Size

var ErrInvalidHMAC = errors.New("Invalid HMAC")
var ErrInvalidData = errors.New("Invalid data")
var ErrWrongPassword = errors.New("Wrong password")
