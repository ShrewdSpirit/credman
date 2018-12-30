package management

import (
	"regexp"
	"sort"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/atotto/clipboard"
)

type SiteData struct {
	ManagementData
	SiteName        string
	NewSiteName     string
	SitePassword    string
	SiteFieldsMap   map[string]string
	SiteFieldsList  []string
	SiteTags        []string
	SiteDeleteTags  []string
	SiteCopyField   bool
	SiteGetTags     bool
	ProfilePassword string
	Profile         *data.Profile
	YesNoPrompt     GetYesNoPromptFunc
	Match           func([]SiteList)
	LogFields       func(SiteFields)
}

type SiteList struct {
	Name       string
	MatchParts [3]string
	Tags       []string
}

type SiteFields struct {
	Fields map[string]string
	Tags   []string
}

func (s SiteData) Add() {
	s.ManagementData.CallStep(SiteStepCheckingExistence)
	if s.Profile.SiteExist(s.SiteName) {
		s.ManagementData.CallStep(SiteStepSiteExists)
		return
	}

	s.ManagementData.CallStep(SiteStepCreating)
	site := make(data.Site)

	s.ManagementData.CallStep(SiteStepSettingFields)
	if len(s.SitePassword) != 0 {
		site["password"] = s.SitePassword
	}
	for field, value := range s.SiteFieldsMap {
		if field == "password" {
			continue
		}
		site[field] = value
	}

	if s.SiteTags != nil && len(s.SiteTags) > 0 {
		s.ManagementData.CallStep(SiteStepAddingTags)
		site.SetTags(s.SiteTags)
	}

	s.ManagementData.CallStep(SiteStepAdding)
	s.Profile.AddSite(s.SiteName, site)

	s.ManagementData.CallStep(ProfileStepSaving)
	if err := s.Profile.Save(s.ProfilePassword); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}

func (s SiteData) Remove() {
	s.ManagementData.CallStep(SiteStepCheckingExistence)
	if !s.Profile.SiteExist(s.SiteName) {
		s.ManagementData.CallStep(SiteStepDoesntExist)
		return
	}

	if !s.YesNoPrompt(SiteStepRemoving) {
		return
	}

	s.ManagementData.CallStep(SiteStepRemoving)
	s.Profile.DeleteSite(s.SiteName)

	s.ManagementData.CallStep(ProfileStepSaving)
	if err := s.Profile.Save(s.ProfilePassword); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}

func (s SiteData) Rename() {
	s.ManagementData.CallStep(SiteStepCheckingExistence)
	if !s.Profile.SiteExist(s.SiteName) {
		s.ManagementData.CallStep(SiteStepDoesntExist)
		return
	}

	s.ManagementData.CallStep(SiteStepRenaming)
	s.Profile.RenameSite(s.SiteName, s.NewSiteName)

	s.ManagementData.CallStep(ProfileStepSaving)
	if err := s.Profile.Save(s.ProfilePassword); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}

func (s SiteData) Set() {
	s.ManagementData.CallStep(SiteStepCheckingExistence)
	if !s.Profile.SiteExist(s.SiteName) {
		s.ManagementData.CallStep(SiteStepDoesntExist)
		return
	}

	site := s.Profile.GetSite(s.SiteName)

	s.ManagementData.CallStep(SiteStepUpdatingFields)
	if len(s.SitePassword) != 0 {
		site["password"] = s.SitePassword
	}
	for field, value := range s.SiteFieldsMap {
		if field == "password" {
			continue
		}
		site[field] = value
	}

	s.ManagementData.CallStep(SiteStepDeletingFields)
	for _, field := range s.SiteFieldsList {
		delete(site, field)
	}

	if s.SiteTags != nil && len(s.SiteTags) > 0 {
		s.ManagementData.CallStep(SiteStepAddingTags)
		site.SetTags(s.SiteTags)
	}

	if s.SiteDeleteTags != nil && len(s.SiteDeleteTags) > 0 {
		s.ManagementData.CallStep(SiteStepRemovingTags)
		site.SetTags(s.SiteTags)
	}

	s.ManagementData.CallStep(ProfileStepSaving)
	if err := s.Profile.Save(s.ProfilePassword); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}

func (s SiteData) List() {
	siteList := make([]SiteList, 0)
	filterTags := s.SiteTags != nil && len(s.SiteTags) > 0
	sortedSiteNames := make([]string, 0, len(s.Profile.Sites))

	for name := range s.Profile.Sites {
		sortedSiteNames = append(sortedSiteNames, name)
	}
	sort.Strings(sortedSiteNames)

	if len(s.SiteName) == 0 {
		for _, name := range sortedSiteNames {
			site := SiteList{
				Name:       name,
				MatchParts: [3]string{name},
			}

			if filterTags {
				if found, tags := s.Profile.Sites[name].HasTags(s.SiteTags); found {
					site.Tags = tags
					siteList = append(siteList, site)
				}
			} else {
				site.Tags = s.Profile.Sites[name].GetTags()
				siteList = append(siteList, site)
			}
		}
	} else {
		s.ManagementData.CallStep(SiteStepRegexCompile)
		rx, err := regexp.Compile(s.SiteName)
		if err != nil {
			s.ManagementData.CallError(SiteStepRegexCompile, err)
			return
		}

		for _, name := range sortedSiteNames {
			if rx.MatchString(name) {
				idx := rx.FindStringIndex(name)

				part1 := name[:idx[0]]
				part2 := name[idx[0]:idx[1]]
				part3 := name[idx[1]:]

				site := SiteList{
					Name:       name,
					MatchParts: [3]string{part1, part2, part3},
				}

				if filterTags {
					if found, tags := s.Profile.Sites[name].HasTags(s.SiteTags); found {
						site.Tags = tags
						siteList = append(siteList, site)
					}
				} else {
					site.Tags = s.Profile.Sites[name].GetTags()
					siteList = append(siteList, site)
				}
			}
		}
	}

	s.Match(siteList)
	s.ManagementData.CallStep(StepDone)
}

func (s SiteData) Get() {
	s.ManagementData.CallStep(SiteStepCheckingExistence)
	if !s.Profile.SiteExist(s.SiteName) {
		s.ManagementData.CallStep(SiteStepDoesntExist)
		return
	}

	site := s.Profile.GetSite(s.SiteName)

	if s.SiteCopyField {
		if len(s.SiteFieldsList) == 0 {
			if err := clipboard.WriteAll(site["password"]); err != nil {
				s.ManagementData.CallError(SiteStepSettingClipboardPassword, err)
				return
			}
			s.ManagementData.CallStep(SiteStepSettingClipboardPassword)
		} else {
			field := s.SiteFieldsList[0]
			_, ok := site[field]
			if !ok || data.IsSpecialField(field) {
				s.ManagementData.CallStep(SiteStepInvalidField)
				return
			}

			if err := clipboard.WriteAll(site[field]); err != nil {
				s.ManagementData.CallError(SiteStepSettingClipboard, err)
				return
			}
			s.ManagementData.CallStep(SiteStepSettingClipboard)
		}
	} else {
		siteFields := SiteFields{
			Fields: make(map[string]string),
		}

		if s.SiteGetTags {
			siteFields.Tags = site.GetTags()
		}

		if len(s.SiteFieldsList) == 0 {
			s.ManagementData.CallStep(SiteStepListingFields)
			for field, value := range site {
				if data.IsSpecialField(field) {
					continue
				}
				siteFields.Fields[field] = value
			}
		} else {
			for _, field := range s.SiteFieldsList {
				if data.IsSpecialField(field) {
					continue
				}

				value, ok := site[field]
				if !ok {
					continue
				}
				siteFields.Fields[field] = value
			}
		}

		s.LogFields(siteFields)
	}

	s.ManagementData.CallStep(StepDone)
}
