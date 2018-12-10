package config

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type UserProfile struct {
	Hash    string
	Content []byte // encrypted binary data of Profile
}

func GetUserProfile(profileDir, name string) (p *UserProfile, err error) {
	profileFilename := path.Join(profileDir, name+".json")
	if _, err = os.Stat(profileFilename); err != nil {
		return
	}

	var data []byte
	if data, err = ioutil.ReadFile(profileFilename); err != nil {
		return
	}

	err = json.Unmarshal(data, p)

	return
}

func NewUserProfile(profileDir, name, password string) error {
	up := &UserProfile{
		Hash: hash(password),
	}
	b, err := json.Marshal(up)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(profileDir, name+".json"), b, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func hash(value string) string {
	h := sha256.Sum256([]byte(value))
	b := []byte(h[:])
	return base64.StdEncoding.EncodeToString(b)
}

func (s *UserProfile) CheckPassword(password string) bool {
	return s.Hash == hash(password)
}
