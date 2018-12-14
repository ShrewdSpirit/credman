package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type SiteList struct {
	List []*Site
}

func newSiteList() SiteList {
	return SiteList{
		List: make([]*Site, 0),
	}
}

type Profile struct {
	Hash  string
	Data  []byte   // encrypted binary data of Profile
	Name  string   `json:"-"`
	Sites SiteList `json:"-"`
}

func GetProfile(profileDir, name string) (*Profile, error) {
	profileFilename := path.Join(profileDir, name+".json")
	if _, err := os.Stat(profileFilename); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(profileFilename)
	if err != nil {
		return nil, err
	}

	p := &Profile{}
	err = json.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}

	p.Sites = newSiteList()
	p.Name = name

	return p, nil
}

func NewProfile(profileDir, name, password string) error {
	p := &Profile{
		Hash: hash(password),
	}
	jsonContent, err := json.Marshal(p)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(profileDir, name+".json"), jsonContent, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (s *Profile) Save(profileDir string) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	profileFilename := path.Join(profileDir, s.Name+".json")
	return ioutil.WriteFile(profileFilename, data, os.ModePerm)
}

func (s *Profile) CheckPassword(password string) bool {
	return s.Hash == hash(password)
}

func (s *Profile) Encrypt(password string) (err error) {
	s.Data, err = json.Marshal(s.Sites)
	if err != nil {
		return
	}

	s.Data, err = Encrypt([]byte(password), s.Data)

	return
}

func (s *Profile) Decrypt(password string) (err error) {
	if len(s.Data) == 0 {
		return
	}

	s.Data, err = Decrypt([]byte(password), s.Data)
	if err != nil {
		return
	}

	err = json.Unmarshal(s.Data, &s.Sites)
	if err != nil {
		return errors.New("Wrong password")
	}

	return
}

func (s *Profile) AddSite(site *Site) {
	s.Sites.List = append(s.Sites.List, site)
}

func (s *Profile) SiteExist(name string) bool {
	for _, site := range s.Sites.List {
		if site.Name == name {
			return true
		}
	}
	return false
}
