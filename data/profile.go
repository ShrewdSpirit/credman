package data

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"time"

	"github.com/ShrewdSpirit/credman/cipher"
)

var ProfileVersion byte = 3

type ProfileMeta struct {
	CreationDate   int64  `json:"c"`
	LastModifyDate int64  `json:"m"`
	Version        byte   `json:"v"`
	Restore        []byte `json:"s"`
	RestoreOrder   []int  `json:"o"`
}

type Profile struct {
	Meta       ProfileMeta `json:"m"`
	SitesBytes []byte      `json:"s"`

	Name  string          `json:"-"`
	Sites map[string]Site `json:"-"`
}

type SiteListResult struct {
	Name       string
	MatchParts [3]string
	Tags       []string
}

func ProfileExists(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	}

	if _, err := os.Stat(path.Join(ProfilesDir, name)); err == nil {
		return true
	}

	return false
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

func LoadProfileRaw(name string) (p *Profile, err error) {
	profileFile := name

	if _, err := os.Stat(name); os.IsNotExist(err) {
		profileFile = path.Join(ProfilesDir, name)
	}

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

	return
}

func LoadProfile(name, password string) (p *Profile, err error) {
	p, err = LoadProfileRaw(name)
	if err != nil {
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

func ListProfiles() ([]string, error) {
	files, err := ioutil.ReadDir(ProfilesDir)
	if err != nil {
		return nil, err
	}

	profileNames := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			profileNames = append(profileNames, file.Name())
		}
	}

	return profileNames, nil
}

func (s *Profile) SaveRaw() (err error) {
	s.Meta.LastModifyDate = time.Now().UnixNano()

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

func (s *Profile) Save(password string) (err error) {
	if s.SitesBytes, err = json.Marshal(s.Sites); err != nil {
		return
	}
	if s.SitesBytes, err = cipher.BlockEncrypt(s.SitesBytes, password); err != nil {
		return
	}

	err = s.SaveRaw()
	return
}

func (s *Profile) AddSite(name string, site Site) {
	s.Sites[name] = site
}

func (s *Profile) SiteExist(name string) bool {
	_, ok := s.Sites[name]
	return ok
}

func (s *Profile) RenameSite(name, newName string) {
	site := s.GetSite(name)
	s.AddSite(newName, site)
	s.DeleteSite(name)
}

func (s *Profile) DeleteSite(name string) {
	delete(s.Sites, name)
}

func (s *Profile) GetSite(name string) Site {
	if s.SiteExist(name) {
		return s.Sites[name]
	}
	return nil
}

func (s *Profile) GetSites(pattern string, tags []string) (result []SiteListResult, err error) {
	result = make([]SiteListResult, 0)
	filterTags := tags != nil && len(tags) > 0
	sortedSiteNames := make([]string, 0, len(s.Sites))

	for siteName := range s.Sites {
		sortedSiteNames = append(sortedSiteNames, siteName)
	}
	sort.Strings(sortedSiteNames)

	if len(pattern) == 0 {
		for _, siteName := range sortedSiteNames {
			siteListItem := SiteListResult{
				Name:       siteName,
				MatchParts: [3]string{siteName},
			}

			if filterTags {
				if found, siteTags := s.Sites[siteName].HasTags(tags); found {
					siteListItem.Tags = siteTags
					result = append(result, siteListItem)
				}
			} else {
				siteListItem.Tags = s.Sites[siteName].GetTags()
				result = append(result, siteListItem)
			}
		}
	} else {
		var rx *regexp.Regexp
		rx, err = regexp.Compile(pattern)
		if err != nil {
			return
		}

		for _, siteName := range sortedSiteNames {
			if rx.MatchString(siteName) {
				locs := rx.FindStringIndex(siteName)

				part1 := siteName[:locs[0]]
				part2 := siteName[locs[0]:locs[1]]
				part3 := siteName[locs[1]:]

				siteListItem := SiteListResult{
					Name:       siteName,
					MatchParts: [3]string{part1, part2, part3},
				}

				if filterTags {
					if found, tags := s.Sites[siteName].HasTags(tags); found {
						siteListItem.Tags = tags
						result = append(result, siteListItem)
					}
				} else {
					siteListItem.Tags = s.Sites[siteName].GetTags()
					result = append(result, siteListItem)
				}
			}
		}
	}

	return
}
