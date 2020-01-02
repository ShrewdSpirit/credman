package cipher

import (
	"bytes"
	"strings"
	"testing"
)

func TestBlock(t *testing.T) {
	data := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam id vehicula augue. Integer dignissim maximus nisl in ultrices. Aenean ultrices ultricies tortor a iaculis.")
	pw := "mzadax7786"

	enc, err := BlockEncrypt(data, pw)
	if err != nil {
		t.Errorf("Encryption error: %s", err)
		return
	}

	metaSize := 64 + 32 + 16 + 32
	if len(enc)-metaSize != len(data) {
		t.Errorf("Encryption wrong: data %d bytes enc %d bytes", len(data), len(enc)-metaSize)
		return
	}

	dec, err := BlockDecrypt(enc, pw)
	if err != nil {
		t.Errorf("Decryption error: %s", err)
	}

	if len(dec) != len(data) {
		t.Errorf("Decryption wrong: data %d bytes dec %d bytes", len(data), len(dec))
	}
}

func TestStream(t *testing.T) {
	data := strings.Repeat("aoijdn1,m2309vckzjv", 100)
	pw := "mamad"

	encout := &bytes.Buffer{}
	if err := StreamEncrypt(strings.NewReader(data), encout, pw); err != nil {
		t.Errorf("Encryption error: %s", err)
		return
	}

	decout := &bytes.Buffer{}
	if err := StreamDecrypt(encout, decout, pw); err != nil {
		t.Errorf("Decryption error: %s", err)
		return
	}

	dec := decout.String()
	if len(dec) != len(data) {
		t.Errorf("Sizes dont match")
		return
	}
	if dec != data {
		t.Errorf("Invalid cipher")
		return
	}
}
