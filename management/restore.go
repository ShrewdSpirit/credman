package management

import (
	"crypto/sha256"

	"github.com/atotto/clipboard"

	"github.com/ShrewdSpirit/credman/cipher"

	"github.com/ShrewdSpirit/credman/data"
)

type RestoreData struct {
	ManagementData
	Answers        []string
	Orders         []int
	Profile        *data.Profile
	PasswordReader GetStringFunc
}

func (s RestoreData) Add() {
	key := s.makeRestoreKey()

	profilePassword := s.PasswordReader(RestoreStepReadingPassword)
	if len(profilePassword) == 0 {
		return
	}

	var err error
	s.ManagementData.CallStep(RestoreStepEncrypting)
	if s.Profile.Meta.Restore, err = cipher.BlockEncrypt([]byte(profilePassword), string(key)); err != nil {
		s.ManagementData.CallError(RestoreStepEncrypting, err)
		return
	}

	s.Profile.Meta.RestoreOrder = s.Orders

	s.ManagementData.CallStep(ProfileStepSaving)
	if err = s.Profile.SaveRaw(); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}

func (s RestoreData) Restore() {
	key := s.makeRestoreKey()

	s.ManagementData.CallStep(RestoreStepDecrypting)
	profilePassword, err := cipher.BlockDecrypt(s.Profile.Meta.Restore, string(key))
	if err != nil {
		s.ManagementData.CallError(RestoreStepDecrypting, err)
		return
	}

	s.ManagementData.CallStep(RestoreStepClipboardPassword)
	if err := clipboard.WriteAll(string(profilePassword)); err != nil {
		s.ManagementData.CallError(RestoreStepClipboardPassword, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}

func (s RestoreData) Remove() {
	s.Profile.Meta.Restore = nil

	s.ManagementData.CallStep(ProfileStepSaving)
	if err := s.Profile.SaveRaw(); err != nil {
		s.ManagementData.CallError(ProfileStepSaving, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}

func (s RestoreData) makeRestoreKey() []byte {
	sha := sha256.New()
	for _, answer := range s.Answers {
		sha.Write([]byte(answer))
	}
	for _, order := range s.Orders {
		sha.Write([]byte{byte(order % 255)})
	}
	sha.Write([]byte{byte(len(s.Answers) % 255)})

	return sha.Sum(nil)
}
