package ezsec

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

func GenerateRSAKeyPair(bits int) (priv *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	priv, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	pub = &priv.PublicKey
	return
}

func RSAPrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

func RSAPublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

func RSABytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error

	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func RSABytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error

	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}

	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Invalid public key")
	}

	return key, nil
}

func RSAEncrypt(shaType ShaType, msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := ShaFunc(shaType)()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func RSADecrypt(shaType ShaType, ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := ShaFunc(shaType)()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func initRSA() {
	js.Module.Get("exports").Set("rsaGenerateKeys", func(bits int) *js.Object {
		r := newResult()

		priv, pub, err := GenerateRSAKeyPair(bits)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		pubB, err := RSAPublicKeyToBytes(pub)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		privB := RSAPrivateKeyToBytes(priv)

		keysObj := js.Global.Get("Object").New()
		keysObj.Set("public", string(pubB))
		keysObj.Set("private", string(privB))
		r.SetValue(keysObj)

		return r.Object
	})

	js.Module.Get("exports").Set("RSAEncrypt", func(shaType ShaType, data interface{}, pub string) *js.Object {
		r := newResult()

		pubKey, err := RSABytesToPublicKey([]byte(pub))
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		data = valueFromInterface(data, valueTypeByteSlice)
		if data == nil {
			r.SetError(ErrInvalidDataType)
			return r.Object
		}
		convertedData := data.([]byte)

		enc, err := RSAEncrypt(ShaType(shaType), convertedData, pubKey)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		r.SetValue(enc)
		return r.Object
	})

	js.Module.Get("exports").Set("RSADecrypt", func(shaType ShaType, data interface{}, priv string) *js.Object {
		r := newResult()

		privKey, err := RSABytesToPrivateKey([]byte(priv))
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		data = valueFromInterface(data, valueTypeByteSlice)
		if data == nil {
			r.SetError(ErrInvalidDataType)
			return r.Object
		}
		convertedData := data.([]byte)

		dec, err := RSADecrypt(ShaType(shaType), convertedData, privKey)
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		r.SetValue(dec)
		return r.Object
	})
}
