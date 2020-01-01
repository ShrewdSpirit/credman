package ezsec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"

	"github.com/gopherjs/gopherjs/js"
)

var ErrWrongPassword = errors.New("Wrong password")

func CFBEncrypt(hmacShaType ShaType, data []byte, key []byte) (encrypted []byte, err error) {
	var aesBlock cipher.Block
	if aesBlock, err = aes.NewCipher(key); err != nil {
		return
	}

	encrypted = make([]byte, len(data))
	iv := RandomBuffer(16)

	cfb := cipher.NewCFBEncrypter(aesBlock, iv)
	cfb.XORKeyStream(encrypted, data)

	hmac := hmac.New(ShaFunc(hmacShaType), key[16:])
	hmac.Write(iv)
	binary.Write(hmac, binary.LittleEndian, len(data))
	hmac.Write(encrypted)

	buf := &bytes.Buffer{}
	buf.Write(iv)
	buf.Write(encrypted)
	buf.Write(hmac.Sum(nil))
	encrypted = buf.Bytes()

	return
}

func CFBDecrypt(hmacShaType ShaType, data []byte, key []byte) (decrypted []byte, err error) {
	buf := bytes.NewReader(data)

	iv := make([]byte, 16)
	if _, err = io.ReadFull(buf, iv); err != nil {
		return
	}

	var aesBlock cipher.Block
	if aesBlock, err = aes.NewCipher(key); err != nil {
		return
	}

	cfb := cipher.NewCFBDecrypter(aesBlock, iv)
	if decrypted, err = ioutil.ReadAll(buf); err != nil {
		return
	}

	macEnd := len(decrypted) - ShaSize(hmacShaType)
	mac := decrypted[macEnd:]
	decrypted = decrypted[:macEnd]

	hmacx := hmac.New(ShaFunc(hmacShaType), key[16:])
	hmacx.Write(iv)
	binary.Write(hmacx, binary.LittleEndian, len(decrypted))
	hmacx.Write(decrypted)

	if !hmac.Equal(mac, hmacx.Sum(nil)) {
		err = ErrWrongPassword
		return
	}

	cfb.XORKeyStream(decrypted, decrypted)

	return
}

func initAES() {
	js.Module.Get("exports").Set("CFBEncrypt", func(hmacShaType int, data interface{}, key interface{}) *js.Object {
		r := newResult()

		key = valueFromInterface(key, valueTypeByteSlice)
		if key == nil {
			r.SetError(ErrInvalidKeyType)
			return r.Object
		}
		convertedKey := key.([]byte)

		data = valueFromInterface(data, valueTypeByteSlice)
		if data == nil {
			r.SetError(ErrInvalidDataType)
			return r.Object
		}
		convertedData := data.([]byte)

		enc, err := CFBEncrypt(ShaType(hmacShaType), convertedData, convertedKey)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		r.SetValue(enc)
		return r.Object
	})

	js.Module.Get("exports").Set("CFBDecrypt", func(hmacShaType int, data interface{}, key interface{}) *js.Object {
		r := newResult()

		key = valueFromInterface(key, valueTypeByteSlice)
		if key == nil {
			r.SetError(ErrInvalidKeyType)
			return r.Object
		}
		convertedKey := key.([]byte)

		data = valueFromInterface(data, valueTypeByteSlice)
		if data == nil {
			r.SetError(ErrInvalidDataType)
			return r.Object
		}
		convertedData := data.([]byte)

		dec, err := CFBDecrypt(ShaType(hmacShaType), convertedData, convertedKey)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		r.SetValue(dec)
		return r.Object
	})
}
