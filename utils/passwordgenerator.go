package utils

import (
	"bytes"
	"math/rand"
	"time"
)

type PasswordCase int
type PasswordMix int

const (
	PasswordCaseBoth  PasswordCase = 0
	PasswordCaseLower PasswordCase = 1
	PasswordCaseUpper PasswordCase = 2
)

const (
	PasswordMixAll    PasswordMix = 0
	PasswordMixLetter PasswordMix = 1
	PasswordMixDigit  PasswordMix = 2
	PasswordMixPunc   PasswordMix = 3
)

var smallLetters = "abcdefghijklmnopqrstuv"
var capLetters = "ABCDEFGHIJKLMNOPQRSTUV"
var digitChars = "0123456789"
var puncChars = " `~!@#$%^&*()_-+=[]{};':\"<>,./?"

func GeneratePassword(plen byte, pcase PasswordCase, pmix ...PasswordMix) string {
	for index, mix := range pmix {
		if mix == PasswordMixAll && len(pmix) > 1 {
			pmix = append(pmix[:index], pmix[index+1:]...)
			break
		}
	}

	mixBuf := bytes.Buffer{}
	for _, mix := range pmix {
		switch mix {
		case PasswordMixAll:
			if pcase == PasswordCaseUpper {
				mixBuf.WriteString(capLetters)
			} else if pcase == PasswordCaseLower {
				mixBuf.WriteString(smallLetters)
			} else {
				mixBuf.WriteString(capLetters)
				mixBuf.WriteString(smallLetters)
			}
			mixBuf.WriteString(digitChars)
			mixBuf.WriteString(puncChars)
		case PasswordMixLetter:
			if pcase == PasswordCaseUpper {
				mixBuf.WriteString(capLetters)
			} else if pcase == PasswordCaseLower {
				mixBuf.WriteString(smallLetters)
			} else {
				mixBuf.WriteString(capLetters)
				mixBuf.WriteString(smallLetters)
			}
		case PasswordMixDigit:
			mixBuf.WriteString(digitChars)
		case PasswordMixPunc:
			mixBuf.WriteString(puncChars)
		}
	}

	srand := rand.NewSource(time.Now().UnixNano())
	mixBytes := mixBuf.Bytes()
	mixLen := int64(len(mixBytes))
	buf := bytes.Buffer{}
	i := byte(0)
	for ; i < plen; i++ {
		buf.WriteByte(mixBytes[srand.Int63()%mixLen])
	}
	return buf.String()
}
