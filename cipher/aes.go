package cipher

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"io/ioutil"
	"reflect"
)

func BlockEncrypt(data []byte, key string) (encrypted []byte, err error) {
	masterHash, hashSalt := HashScrypt(key, ScryptDifficultyNorm, 128)
	keyAes := masterHash[:32]
	keyHmac := masterHash[32:64]
	checkHash := masterHash[64:]
	iv := GenerateIV()

	var aesBlock cipher.Block
	if aesBlock, err = aes.NewCipher(keyAes); err != nil {
		return
	}

	encrypted = make([]byte, len(data))

	cfb := cipher.NewCFBEncrypter(aesBlock, iv)
	cfb.XORKeyStream(encrypted, data)

	hmacX := hmac.New(sha256.New, keyHmac)
	hmacX.Write(iv)
	binary.Write(hmacX, binary.LittleEndian, uint64(len(data)))
	hmacX.Write(encrypted)

	buf := &bytes.Buffer{}
	buf.Write(checkHash)
	buf.Write(hashSalt)
	buf.Write(iv)
	buf.Write(encrypted)
	buf.Write(hmacX.Sum(nil))
	encrypted = buf.Bytes()

	return
}

func BlockDecrypt(data []byte, key string) (decrypted []byte, err error) {
	buf := bytes.NewReader(data)

	checkHash := make([]byte, 64)
	if _, err = io.ReadFull(buf, checkHash); err != nil {
		return
	}

	hashSalt := make([]byte, 32)
	if _, err = io.ReadFull(buf, hashSalt); err != nil {
		return
	}

	masterHash := HashScryptSalt(key, ScryptDifficultyNorm, 128, hashSalt)
	keyAes := masterHash[:32]
	keyHmac := masterHash[32:64]
	keyCheckHash := masterHash[64:]

	if !reflect.DeepEqual(keyCheckHash, checkHash) {
		return nil, ErrWrongPassword
	}

	iv := make([]byte, 16)
	if _, err = io.ReadFull(buf, iv); err != nil {
		return
	}

	var aesBlock cipher.Block
	if aesBlock, err = aes.NewCipher(keyAes); err != nil {
		return
	}

	cfb := cipher.NewCFBDecrypter(aesBlock, iv)
	if decrypted, err = ioutil.ReadAll(buf); err != nil {
		return
	}

	decSize := len(decrypted)
	mac := decrypted[decSize-hmacSize : decSize]
	decrypted = decrypted[:decSize-hmacSize]

	hmacX := hmac.New(sha256.New, keyHmac)
	hmacX.Write(iv)
	binary.Write(hmacX, binary.LittleEndian, uint64(len(decrypted)))
	hmacX.Write(decrypted)

	if !hmac.Equal(mac, hmacX.Sum(nil)) {
		err = ErrInvalidHMAC
		return
	}

	cfb.XORKeyStream(decrypted, decrypted)

	return
}

func StreamEncrypt(in io.Reader, out io.Writer, key string) error {
	masterHash, hashSalt := HashScrypt(key, ScryptDifficultyNorm, 128)
	keyAes := masterHash[:32]
	keyHmac := masterHash[32:64]
	checkHash := masterHash[64:]
	iv := GenerateIV()

	aesBlock, err := aes.NewCipher(keyAes)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aesBlock, iv)
	hmac := hmac.New(sha256.New, keyHmac)

	out.Write(checkHash)
	out.Write(hashSalt)

	w := io.MultiWriter(out, hmac)
	w.Write(iv)

	buf := make([]byte, bufSize)
	for {
		var bufSz int
		bufSz, err = in.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			break
		}

		if bufSz != 0 {
			outBuf := make([]byte, bufSz)
			ctr.XORKeyStream(outBuf, buf[:bufSz])
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

	masterHash := HashScryptSalt(key, ScryptDifficultyNorm, 128, hashSalt)
	keyAes := masterHash[:32]
	keyHmac := masterHash[32:64]
	keyCheckHash := masterHash[64:]

	if !reflect.DeepEqual(keyCheckHash, checkHash) {
		return ErrWrongPassword
	}

	iv := make([]byte, 16)
	if _, err := io.ReadFull(in, iv); err != nil {
		return err
	}

	aesBlock, err := aes.NewCipher(keyAes)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aesBlock, iv)
	hmacX := hmac.New(sha256.New, keyHmac)
	hmacX.Write(iv)

	mac := make([]byte, hmacSize)
	buf := bufio.NewReaderSize(in, bufSize)
	for {
		readBytes, err := buf.Peek(bufSize)
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			nBytesLeft := buf.Buffered()
			if nBytesLeft < hmacSize {
				return ErrInvalidData
			}

			copy(mac, readBytes[nBytesLeft-hmacSize:nBytesLeft])
			if nBytesLeft == hmacSize {
				break
			}
		}

		limit := len(readBytes) - hmacSize
		limitedBytes := readBytes[:limit]
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
