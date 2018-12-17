package data

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/ShrewdSpirit/credman/utility"
)

var ProfileVersion string
var ProfilesDir string

func init() {
	ProfileVersion = utility.Hash("1.2")
}

type ProfileMeta struct {
	Remote         *Remote `json:"r"`
	CreationDate   int64   `json:"c"`
	LastModifyDate int64   `json:"m"`
}

type Profile struct {
	Hash    string `json:"h"`
	Meta    []byte `json:"m"`
	Version string `json:"v"`

	Name        string      `json:"-"`
	ProfileMeta ProfileMeta `json:"-"`
	Sites       []*Site     `json:"-"`
}

func ProfileExists(name string) bool {
	_, err := os.Stat(path.Join(ProfilesDir, name))
	if err != nil {
		return false
	}
	return true
}

func NewProfile(name, password string) *Profile {
	return &Profile{
		Hash:    utility.Hash(password),
		Name:    name,
		Version: ProfileVersion,
		Sites:   make([]*Site, 0),
	}
}

func LoadProfile(name, password string) (p *Profile, err error) {
	key := []byte(password)
	profileDir := path.Join(ProfilesDir, name)

	var profileBytes []byte
	if profileBytes, err = ioutil.ReadFile(path.Join(profileDir, "profile.json")); err != nil {
		return
	}

	p = NewProfile(name, password)
	if err = json.Unmarshal(profileBytes, p); err != nil {
		return
	}

	if p.Version != ProfileVersion {
		err = errors.New("Unsupported profile version")
		return
	}

	if !p.CheckPassword(password) {
		err = errors.New("Wrong password")
		return
	}

	if p.Meta, err = utility.Decrypt(key, p.Meta); err != nil {
		return
	}
	if err = json.Unmarshal(p.Meta, &p.ProfileMeta); err != nil {
		return
	}

	var sitesBytes []byte
	if sitesBytes, err = ioutil.ReadFile(path.Join(profileDir, "sites")); err != nil {
		return
	}
	if sitesBytes, err = utility.Decrypt(key, sitesBytes); err != nil {
		return
	}
	if err = json.Unmarshal(sitesBytes, &p.Sites); err != nil {
		return
	}

	return
}

func (s *Profile) Save(password string) (err error) {
	key := []byte(password)

	s.ProfileMeta.LastModifyDate = time.Now().UnixNano()

	if s.Meta, err = json.Marshal(s.ProfileMeta); err != nil {
		return
	}
	if s.Meta, err = utility.Encrypt(key, s.Meta); err != nil {
		return
	}

	var profileBytes []byte
	if profileBytes, err = json.Marshal(s); err != nil {
		return
	}

	var sitesBytes []byte
	if sitesBytes, err = json.Marshal(s.Sites); err != nil {
		return
	}
	if sitesBytes, err = utility.Encrypt(key, sitesBytes); err != nil {
		return
	}

	profileDir := path.Join(ProfilesDir, s.Name)
	if err = ioutil.WriteFile(path.Join(profileDir, "profile.json"), profileBytes, os.ModePerm); err != nil {
		return
	}
	if err = ioutil.WriteFile(path.Join(profileDir, "sites"), sitesBytes, os.ModePerm); err != nil {
		return
	}

	return
}

func (s *Profile) CheckPassword(password string) bool {
	return utility.Hash(password) == s.Hash
}

func (s *Profile) AddSite(site *Site) {
	s.Sites = append(s.Sites, site)
}

func (s *Profile) SiteExist(name string) bool {
	for _, site := range s.Sites {
		if site.Name == name {
			return true
		}
	}
	return false
}

func (s *Profile) GetSite(name string) *Site {
	for _, site := range s.Sites {
		if site.Name == name {
			return site
		}
	}
	return nil
}

func (s *Profile) DeleteSite(name string) {
	siteIndex := -1
	for i, site := range s.Sites {
		if site.Name == name {
			siteIndex = i
			break
		}
	}
	if siteIndex >= 0 {
		s.Sites = append(s.Sites[:siteIndex], s.Sites[siteIndex+1:]...)
	}
}
