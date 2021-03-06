package data

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	"github.com/ShrewdSpirit/credman/utils/vars"
)

type Site map[string]string // fields

const (
	SpecialFieldTags          = "$$$TAGS"
	SpecialFieldFileData      = "$$$FBYTES"
	SpecialFieldFileStoreType = "$$$FSTO"

	FileFieldName     string = "name"
	FileFieldAbsolute string = "path"
	FileFieldUUID     string = "uuid"
	FileFieldAddDate  string = "added"
	FileFieldLastDec  string = "decrypted"
	FileFieldUpdate   string = "updated"
	FileFieldSize     string = "size"
)

type FileStoreType string

const (
	FileStoreTypeBase64 FileStoreType = "base64" // it's the default
	FileStoreTypeHex    FileStoreType = "hex"
)

func IsSpecialField(name string) bool {
	return strings.HasPrefix(name, "$$$")
}

func NewSite(name, password string, fields map[string]string, tags []string) (site Site) {
	site = make(Site)

	if len(password) != 0 {
		site["password"] = password
	}

	for field, value := range fields {
		if field == "password" {
			continue
		}

		site[strings.ToLower(field)] = value
	}

	if tags != nil && len(tags) > 0 {
		site.AddTags(tags)
	}

	return
}

func (s Site) GetFields(filterFields []string) (fields map[string]string) {
	fields = make(map[string]string)

	if filterFields == nil || len(filterFields) == 0 {
		for field, value := range s {
			if IsSpecialField(field) {
				continue
			}
			fields[field] = value
		}
	} else {
		for _, field := range filterFields {
			if IsSpecialField(field) {
				continue
			}

			value, ok := s[field]
			if !ok {
				continue
			}
			fields[field] = value
		}
	}

	return
}

func (s Site) GetTags() (tags []string) {
	tagsList, ok := s[SpecialFieldTags]
	tags = make([]string, 0)
	if !ok {
		return
	}

	json.Unmarshal([]byte(tagsList), &tags)

	return
}

func (s Site) AddTags(tags []string) {
	for _, siteTag := range s.GetTags() {
		found := false
		for _, tag := range tags {
			if tag == siteTag {
				found = true
			}
		}
		if !found {
			tags = append(tags, strings.ToLower(siteTag))
		}
	}

	s.SetTags(tags)
}

func (s Site) SetTags(tags []string) {
	for i, t := range tags {
		tags[i] = strings.ToLower(t)
	}

	tagsBytes, err := json.Marshal(tags)
	if err != nil {
		return
	}

	s[SpecialFieldTags] = string(tagsBytes)
}

func (s Site) HasTag(tag string) bool {
	siteTags := s.GetTags()
	tag = strings.ToLower(tag)

	for _, siteTag := range siteTags {
		if tag == siteTag {
			return true
		}
	}

	return false
}

func (s Site) HasTags(tags []string) (found bool, foundTags []string) {
	foundTags = make([]string, 0)
	siteTags := s.GetTags()

	for i, t := range tags {
		tags[i] = strings.ToLower(t)
	}

	for _, tag := range siteTags {
		for _, askedTag := range tags {
			if tag == askedTag {
				foundTags = append(foundTags, tag)
			}
		}
	}

	found = len(foundTags) > 0

	return
}

func (s Site) RemoveTags(tags []string) {
	siteTags := s.GetTags()
	newTags := make([]string, 0)

	for i, t := range tags {
		tags[i] = strings.ToLower(t)
	}

	for _, siteTag := range siteTags {
		found := false
		for _, tag := range tags {
			if tag == siteTag {
				found = true
			}
		}
		if !found {
			newTags = append(newTags, siteTag)
		}
	}

	s.SetTags(newTags)
}

func NewSiteFile(fileName, absPath, uuid string) (site Site) {
	site = make(Site)
	site[FileFieldName] = fileName
	site[FileFieldAbsolute] = absPath
	site[FileFieldUUID] = uuid
	site[FileFieldAddDate] = time.Now().Local().Format(vars.TimeStringFormat)
	site.AddTags([]string{"file"})
	return
}

func (s Site) IsFile() bool {
	_, ok := s[SpecialFieldFileData]
	return ok
}

func (s Site) GetFileBytes(fs FileStoreType) (result []byte, err error) {
	switch fs {
	case FileStoreTypeBase64:
		result, err = base64.URLEncoding.DecodeString(s[SpecialFieldFileData])
		if err != nil {
			return
		}
	case FileStoreTypeHex:
		result, err = hex.DecodeString(s[SpecialFieldFileData])
		if err != nil {
			return
		}
	}

	return
}

func (s Site) WriteFileBytes(data []byte, fs FileStoreType) {
	switch fs {
	case FileStoreTypeBase64:
		s[SpecialFieldFileData] = base64.URLEncoding.EncodeToString(data)
	case FileStoreTypeHex:
		s[SpecialFieldFileData] = hex.EncodeToString(data)
	}
}

func (s Site) FileStoreType() FileStoreType {
	if fs, ok := s[SpecialFieldFileStoreType]; ok {
		return FileStoreType(fs)
	}

	return FileStoreTypeBase64
}
