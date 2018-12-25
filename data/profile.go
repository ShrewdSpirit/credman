package data

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/ShrewdSpirit/credman/cipher"
)

var ProfileVersion byte = 2
var ProfilesDir string

type ProfileMeta struct {
	Remote         []byte `json:"r"`
	CreationDate   int64  `json:"c"`
	LastModifyDate int64  `json:"m"`
	Version        byte   `json:"v"`
	Restore        []byte `json:"s"`
}

type Site map[string]string // fields

type Profile struct {
	Meta       ProfileMeta `json:"m"`
	SitesBytes []byte      `json:"s"`

	Name  string          `json:"-"`
	Sites map[string]Site `json:"-"`
}

func ProfileExists(name string) bool {
	_, err := os.Stat(path.Join(ProfilesDir, name))
	if err != nil {
		return false
	}
	return true
}

func NewProfile(name string) *Profile {
	return &Profile{
		Name:  name,
		Sites: make(map[string]Site),
		Meta: ProfileMeta{
			CreationDate:   time.Now().UnixNano(),
			LastModifyDate: time.Now().UnixNano(),
			Version:        ProfileVersion,
		},
	}
}

func LoadProfile(name, password string) (p *Profile, err error) {
	profileFile := path.Join(ProfilesDir, name)

	var profileBytes []byte
	if profileBytes, err = ioutil.ReadFile(profileFile); err != nil {
		return
	}

	p = NewProfile(name)
	if err = json.Unmarshal(profileBytes, p); err != nil {
		return
	}

	if p.Meta.Version != ProfileVersion {
		err = errors.New("Unsupported profile version")
		return
	}

	if p.SitesBytes, err = cipher.BlockDecrypt(p.SitesBytes, password); err != nil {
		return
	}
	if err = json.Unmarshal(p.SitesBytes, &p.Sites); err != nil {
		return
	}

	return
}

func (s *Profile) Save(password string) (err error) {
	s.Meta.LastModifyDate = time.Now().UnixNano()

	if s.SitesBytes, err = json.Marshal(s.Sites); err != nil {
		return
	}
	if s.SitesBytes, err = cipher.BlockEncrypt(s.SitesBytes, password); err != nil {
		return
	}

	var profileBytes []byte
	if profileBytes, err = json.Marshal(s); err != nil {
		return
	}

	profileName := path.Join(ProfilesDir, s.Name)
	if err = ioutil.WriteFile(profileName, profileBytes, os.ModePerm); err != nil {
		return
	}

	return
}

func (s *Profile) AddSite(name string, site Site) {
	s.Sites[name] = site
}

func (s *Profile) SiteExist(name string) bool {
	_, ok := s.Sites[name]
	return ok
}

func (s *Profile) GetSite(name string) Site {
	if s.SiteExist(name) {
		return s.Sites[name]
	}
	return nil
}

func (s *Profile) RenameSite(name, newName string) {
	site := s.GetSite(name)
	s.AddSite(newName, site)
	s.DeleteSite(name)
}

func (s *Profile) DeleteSite(name string) {
	delete(s.Sites, name)
}
