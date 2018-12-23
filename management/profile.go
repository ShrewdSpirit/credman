package management

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ShrewdSpirit/credman/data"
)

const (
	ProfileStepCheckingExistence ManagementStep = 0
	ProfileStepProfileExists     ManagementStep = 1
	ProfileStepReadingPassword   ManagementStep = 2
	ProfileStepCreating          ManagementStep = 3
	ProfileStepSaving            ManagementStep = 4
	ProfileStepDone              ManagementStep = 5
	ProfileStepDefaultChanged    ManagementStep = 6
	ProfileStepDoesntExist       ManagementStep = 7
	ProfileStepRemovePrompt      ManagementStep = 8
	ProfileStepRemoving          ManagementStep = 9
	ProfileStepRenaming          ManagementStep = 10
	ProfileStepNewPassword       ManagementStep = 11
	ProfileStepLoadingProfile    ManagementStep = 12
	ProfileStepReadingProfiles   ManagementStep = 13
)

type ProfileData struct {
	ManagementData
	ProfileName    string
	NewProfileName string
	PasswordReader GetStringFunc
	YesNoPrompt    GetYesNoPromptFunc
	LogList        LogStringFunc
}

func (s ProfileData) Add() {
	s.ManagementData.CallStep(ProfileStepCheckingExistence)
	if data.ProfileExists(s.ProfileName) {
		s.ManagementData.CallStep(ProfileStepProfileExists)
		return
	}

	s.ManagementData.CallStep(ProfileStepReadingPassword)
	password := s.PasswordReader(ProfileStepReadingPassword)
	if len(password) == 0 {
		return
	}

	s.ManagementData.CallStep(ProfileStepCreating)
	profile := data.NewProfile(s.ProfileName)

	s.ManagementData.CallStep(ProfileStepSaving)
	if err := profile.Save(password); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}
	s.ManagementData.CallStep(ProfileStepDone)

	if data.Config.DefaultProfile == "" {
		data.Config.DefaultProfile = s.ProfileName
		s.ManagementData.CallStep(ProfileStepDefaultChanged)
	}
}

func (s ProfileData) Remove() {
	s.ManagementData.CallStep(ProfileStepCheckingExistence)
	if !data.ProfileExists(s.ProfileName) {
		s.ManagementData.CallStep(ProfileStepDoesntExist)
		return
	}

	if !s.YesNoPrompt(ProfileStepRemovePrompt) {
		return
	}

	s.ManagementData.CallStep(ProfileStepRemoving)
	if err := os.RemoveAll(path.Join(data.ProfilesDir, s.ProfileName)); err != nil {
		s.ManagementData.CallError(ProfileStepRemoving, err)
		return
	}
	s.ManagementData.CallStep(ProfileStepDone)

	if data.Config.DefaultProfile == s.ProfileName {
		data.Config.DefaultProfile = ""
		s.ManagementData.CallStep(ProfileStepDefaultChanged)
	}
}

func (s ProfileData) Rename() {
	s.ManagementData.CallStep(ProfileStepCheckingExistence)
	if !data.ProfileExists(s.ProfileName) {
		s.ManagementData.CallStep(ProfileStepDoesntExist)
		return
	}

	if data.ProfileExists(s.NewProfileName) {
		s.ManagementData.CallStep(ProfileStepProfileExists)
		return
	}

	s.ManagementData.CallStep(ProfileStepRenaming)
	ppath := path.Join(data.ProfilesDir, s.ProfileName)
	npath := path.Join(data.ProfilesDir, s.NewProfileName)
	if err := os.Rename(ppath, npath); err != nil {
		s.ManagementData.CallError(ProfileStepRenaming, err)
		return
	}
	s.ManagementData.CallStep(ProfileStepDone)

	if data.Config.DefaultProfile == s.ProfileName {
		data.Config.DefaultProfile = s.NewProfileName
		s.ManagementData.CallStep(ProfileStepDefaultChanged)
	}
}

func (s ProfileData) Passwd() {
	s.ManagementData.CallStep(ProfileStepCheckingExistence)
	if !data.ProfileExists(s.ProfileName) {
		s.ManagementData.CallStep(ProfileStepDoesntExist)
		return
	}

	s.ManagementData.CallStep(ProfileStepReadingPassword)
	password := s.PasswordReader(ProfileStepReadingPassword)
	if len(password) == 0 {
		return
	}

	s.ManagementData.CallStep(ProfileStepLoadingProfile)
	profile, err := data.LoadProfile(s.ProfileName, password)
	if err != nil {
		s.ManagementData.CallError(ProfileStepLoadingProfile, err)
		return
	}

	s.ManagementData.CallStep(ProfileStepNewPassword)
	newPassword := s.PasswordReader(ProfileStepNewPassword)
	if len(newPassword) == 0 {
		return
	}

	s.ManagementData.CallStep(ProfileStepSaving)
	if err := profile.Save(newPassword); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}
	s.ManagementData.CallStep(ProfileStepDone)
}

func (s ProfileData) List() {
	s.ManagementData.CallStep(ProfileStepReadingProfiles)
	profiles, err := ioutil.ReadDir(data.ProfilesDir)
	if err != nil {
		s.ManagementData.CallError(ProfileStepReadingProfiles, err)
		return
	}

	for _, profile := range profiles {
		if !profile.IsDir() {
			s.LogList(profile.Name())
		}
	}
	s.ManagementData.CallStep(ProfileStepDone)
}
