package cipher

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
)

const bufSize = 16 * 1024
const hmacSize = sha256.Size

var ErrInvalidHMAC = errors.New("Invalid HMAC")
var ErrInvalidData = errors.New("Invalid data")
var ErrWrongPassword = errors.New("Wrong password")

func BlockEncrypt(data []byte, key string) (encBuffer []byte, err error) {
	aesKey, hmacKey, check, salt := BuildKeys(key, ScryptDifficultyNorm, nil)
	iv := GenerateIV()

	var aesBlock cipher.Block
	if aesBlock, err = aes.NewCipher(aesKey); err != nil {
		return
	}

	dataLen := len(data)
	encBuffer = make([]byte, dataLen)
	cfb := cipher.NewCFBEncrypter(aesBlock, iv)
	cfb.XORKeyStream(encBuffer, data)

	hmacHash := BuildHmac(uint64(dataLen), encBuffer, hmacKey, iv).Sum(nil)

	buf := &bytes.Buffer{}
	buf.Write(check)
	buf.Write(salt)
	buf.Write(iv)
	buf.Write(encBuffer)
	buf.Write(hmacHash)

	encBuffer = buf.Bytes()

	return
}

func BlockDecrypt(data []byte, key string) (decBuffer []byte, err error) {
	buf := bytes.NewReader(data)

	bufCheck := make([]byte, 64)
	if _, err = io.ReadFull(buf, bufCheck); err != nil {
		return
	}

	salt := make([]byte, 32)
	if _, err = io.ReadFull(buf, salt); err != nil {
		return
	}

	aesKey, hmacKey, check, _ := BuildKeys(key, ScryptDifficultyNorm, salt)
	if !reflect.DeepEqual(bufCheck, check) {
		err = ErrWrongPassword
		return
	}

	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(buf, iv); err != nil {
		return
	}

	var aesBlock cipher.Block
	if aesBlock, err = aes.NewCipher(aesKey); err != nil {
		return
	}

	cfb := cipher.NewCFBDecrypter(aesBlock, iv)
	if decBuffer, err = ioutil.ReadAll(buf); err != nil {
		return
	}

	decLen := len(decBuffer)
	hmacHash := decBuffer[decLen-hmacSize : decLen]
	decBuffer = decBuffer[:decLen-hmacSize]

	newHmac := BuildHmac(uint64(len(decBuffer)), decBuffer, hmacKey, iv)
	if !hmac.Equal(hmacHash, newHmac.Sum(nil)) {
		err = ErrInvalidHMAC
		return
	}

	cfb.XORKeyStream(decBuffer, decBuffer)

	return
}

func StreamEncrypt(in io.Reader, out io.Writer, key string) error {
	aesKey, hmacKey, check, salt := BuildKeys(key, ScryptDifficultyNorm, nil)
	iv := GenerateIV()

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aesBlock, iv)
	hmac := hmac.New(sha256.New, hmacKey)

	out.Write(check)
	out.Write(salt)

	w := io.MultiWriter(out, hmac)
	w.Write(iv)

	buf := make([]byte, bufSize)
	for {
		var numBytesRead int
		numBytesRead, err = in.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			break
		}

		if numBytesRead != 0 {
			outBuf := make([]byte, numBytesRead)
			ctr.XORKeyStream(outBuf, buf[:numBytesRead])
			w.Write(outBuf)
		}
	}

	out.Write(hmac.Sum(nil))

	return nil
}

func StreamDecrypt(in io.Reader, out io.Writer, key string) error {
	checkHash := make([]byte, 64)
	if _, err := io.ReadFull(in, checkHash); err != nil {
		return err
	}

	hashSalt := make([]byte, 32)
	if _, err := io.ReadFull(in, hashSalt); err != nil {
		return err
	}

	aesKey, hmacKey, keyCheck, _ := BuildKeys(key, ScryptDifficultyNorm, hashSalt)

	if !reflect.DeepEqual(keyCheck, checkHash) {
		return ErrWrongPassword
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(in, iv); err != nil {
		return err
	}

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aesBlock, iv)
	hmacX := hmac.New(sha256.New, hmacKey)
	hmacX.Write(iv)

	mac := make([]byte, hmacSize)
	buf := bufio.NewReaderSize(in, bufSize)
	for {
		numBytesRead, err := buf.Peek(bufSize)
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			numBytesLeft := buf.Buffered()
			if numBytesLeft < hmacSize {
				return ErrInvalidData
			}

			copy(mac, numBytesRead[numBytesLeft-hmacSize:numBytesLeft])
			if numBytesLeft == hmacSize {
				break
			}
		}

		limit := len(numBytesRead) - hmacSize
		limitedBytes := numBytesRead[:limit]
		outbuf := make([]byte, limit)

		hmacX.Write(limitedBytes)
		buf.Read(limitedBytes)
		ctr.XORKeyStream(outbuf, limitedBytes)
		out.Write(outbuf)

		if err == io.EOF {
			break
		}
	}

	if !hmac.Equal(mac, hmacX.Sum(nil)) {
		return ErrInvalidHMAC
	}

	return nil
}
