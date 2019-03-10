package data

import (
	"encoding/json"
	"strings"
)

type Site map[string]string // fields

const SpecialFieldTags = "$$$TAGS"
const SpecialFieldFileData = "$$$FBYTES" // dont store em here

func IsSpecialField(name string) bool {
	return strings.HasPrefix(name, "$$$")
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
