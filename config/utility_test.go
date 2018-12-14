package config

import "testing"

func TestAES(t *testing.T) {
	data := "Hello worldddd! w"
	key := "12345678901234"

	enc, err := Encrypt([]byte(key), []byte(data))
	if err != nil {
		t.Errorf("Encryption failed: %s", err)
		return
	}

	dec, err := Decrypt([]byte(key), enc)
	if err != nil {
		t.Errorf("Decrption failed: %s", err)
		return
	}

	if string(dec) != data {
		t.Error("Wrong decryption")
	}
}
