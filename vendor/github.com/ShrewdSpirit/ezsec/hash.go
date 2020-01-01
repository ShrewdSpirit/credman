package ezsec

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"hash"
	"hash/fnv"

	"github.com/gopherjs/gopherjs/js"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

type ShaType byte

const (
	ShaTypeSha1 ShaType = iota
	ShaTypeSha256
	ShaTypeSha512
)

func ShaFunc(shaType ShaType) func() hash.Hash {
	switch shaType {
	case ShaTypeSha1:
		return sha1.New
	case ShaTypeSha256:
		return sha256.New
	case ShaTypeSha512:
		return sha512.New
	default:
		return nil
	}
}

func ShaSize(shaType ShaType) int {
	switch shaType {
	case ShaTypeSha1:
		return sha1.Size
	case ShaTypeSha256:
		return sha256.Size
	case ShaTypeSha512:
		return sha512.Size
	default:
		return 0
	}
}

func Sha(shaType ShaType, value string) []byte {
	hash := ShaFunc(shaType)()
	hash.Write([]byte(value))
	return hash.Sum(nil)
}

func Bcrypt(key []byte, cost int) (hash []byte, err error) {
	hash, err = bcrypt.GenerateFromPassword(key, cost)
	return
}

const (
	ScryptCostLow    = 8
	ScryptCostNormal = 16
	ScyrptCostHigh   = 24
)

func Scrypt(value []byte, cost int, hashLength int, salt []byte) (hash []byte, o_salt []byte, err error) {
	if salt == nil || len(salt) == 0 {
		salt = RandomBuffer(32)
	}

	hash, err = scrypt.Key(value, salt, 65536, cost, 1, hashLength)
	o_salt = salt
	return
}

func Fnv1a32(value string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(value))
	return binary.LittleEndian.Uint32(hash.Sum(nil))
}

func Fnv1a64(value string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(value))
	return binary.LittleEndian.Uint64(hash.Sum(nil))
}

func Base64Encode(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

func Base64Decode(value string) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(value)
	return b, err
}

func initHash() {
	scryptCostObj := js.Global.Get("Object").New()
	scryptCostObj.Set("low", ScryptCostLow)
	scryptCostObj.Set("normal", ScryptCostNormal)
	scryptCostObj.Set("high", ScyrptCostHigh)
	js.Module.Get("exports").Set("ScryptCost", scryptCostObj)

	shaTypeObj := js.Global.Get("Object").New()
	shaTypeObj.Set("sha1", ShaTypeSha1)
	shaTypeObj.Set("sha256", ShaTypeSha256)
	shaTypeObj.Set("sha512", ShaTypeSha512)
	js.Module.Get("exports").Set("ShaType", shaTypeObj)

	js.Module.Get("exports").Set("bcrypt", func(value interface{}, cost int) *js.Object {
		r := newResult()

		value = valueFromInterface(value, valueTypeByteSlice)
		if value == nil {
			r.SetError(ErrInvalidDataType)
			return r.Object
		}
		convertedValue := value.([]byte)

		hash, err := Bcrypt(convertedValue, cost)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		r.SetValue(hash)
		return r.Object
	})

	js.Module.Get("exports").Set("scrypt", func(value interface{}, cost int, hashLength int, a_salt interface{}) *js.Object {
		r := newResult()

		value = valueFromInterface(value, valueTypeByteSlice)
		if value == nil {
			r.SetError(ErrInvalidDataType)
			return r.Object
		}
		convertedValue := value.([]byte)

		var salt []byte
		switch a_salt.(type) {
		case []byte:
			salt = a_salt.([]byte)
		}

		hash, salt, err := Scrypt(convertedValue, cost, hashLength, salt)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		scryptObj := js.Global.Get("Object").New()
		scryptObj.Set("hash", hash)
		scryptObj.Set("salt", salt)
		r.SetValue(scryptObj)

		return r.Object
	})

	js.Module.Get("exports").Set("sha", func(shaType int, value string) []byte {
		return Sha(ShaType(shaType), value)
	})

	js.Module.Get("exports").Set("fnv1a32", func(value interface{}) uint32 {
		switch value.(type) {
		case string:
			return Fnv1a32(value.(string))
		case []byte:
			return Fnv1a32(string(value.([]byte)))
		}
		return 0
	})

	js.Module.Get("exports").Set("fnv1a64", func(value interface{}) uint64 {
		switch value.(type) {
		case string:
			return Fnv1a64(value.(string))
		case []byte:
			return Fnv1a64(string(value.([]byte)))
		}
		return 0
	})

	js.Module.Get("exports").Set("base64Encode", func(value interface{}) string {
		switch value.(type) {
		case string:
			return Base64Encode([]byte(value.(string)))
		case []byte:
			return Base64Encode(value.([]byte))
		}
		return ""
	})

	js.Module.Get("exports").Set("base64Decode", func(value string) *js.Object {
		r := newResult()

		result, err := Base64Decode(value)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		r.SetValue(result)
		return r.Object
	})
}
