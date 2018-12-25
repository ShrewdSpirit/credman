package management

import (
	"regexp"

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
	SiteCopyField   bool
	ProfilePassword string
	Profile         *data.Profile
	YesNoPrompt     GetYesNoPromptFunc
	Match           LogMatch
	LogFields       LogString2Func
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

	s.ManagementData.CallStep(ProfileStepSaving)
	if err := s.Profile.Save(s.ProfilePassword); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}

func (s SiteData) List() {
	if len(s.SiteName) == 0 {
		for name, _ := range s.Profile.Sites {
			s.Match(name, name, "", "")
		}
	} else {
		s.ManagementData.CallStep(SiteStepRegexCompile)
		rx, err := regexp.Compile(s.SiteName)
		if err != nil {
			s.ManagementData.CallError(SiteStepRegexCompile, err)
			return
		}

		for name, _ := range s.Profile.Sites {
			if rx.MatchString(name) {
				idx := rx.FindStringIndex(name)
				part1 := name[:idx[0]]
				part2 := name[idx[0]:idx[1]]
				part3 := name[idx[1]:]
				s.Match(name, part1, part2, part3)
			}
		}
	}

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
			if !ok {
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
		if len(s.SiteFieldsList) == 0 {
			s.ManagementData.CallStep(SiteStepListingFields)
			for field, value := range site {
				s.LogFields(field, value)
			}
		} else {
			for _, field := range s.SiteFieldsList {
				value, ok := site[field]
				if !ok {
					continue
				}
				s.LogFields(field, value)
			}
		}
	}

	s.ManagementData.CallStep(StepDone)
}
