package management

import (
	"os"
	"strings"

	"github.com/ShrewdSpirit/credman/cipher"
)

type FileData struct {
	ManagementData
	InputFilename  string
	OutputFilename *string
	DeleteInput    bool
	PasswordReader GetStringFunc
}

func (s FileData) Encrypt() {
	s.ManagementData.CallStep(FileStepCheckingExistence)
	_, err := os.Stat(s.InputFilename)
	if err != nil {
		s.ManagementData.CallError(FileStepCheckingExistence, err)
		return
	}

	if len(*s.OutputFilename) == 0 {
		*s.OutputFilename = s.InputFilename + ".enc"
	}

	password := s.PasswordReader(FileStepReadingPassword)
	if len(password) == 0 {
		return
	}

	s.ManagementData.CallStep(FileStepOpeningInput)
	inputFile, err := os.Open(s.InputFilename)
	if err != nil {
		s.ManagementData.CallError(FileStepOpeningInput, err)
		return
	}

	s.ManagementData.CallStep(FileStepCreatingOutput)
	outputFile, err := os.Create(*s.OutputFilename)
	if err != nil {
		s.ManagementData.CallError(FileStepCreatingOutput, err)
		return
	}

	s.ManagementData.CallStep(FileStepEncrypting)
	if err = cipher.StreamEncrypt(inputFile, outputFile, password); err != nil {
		os.Remove(*s.OutputFilename)
		s.ManagementData.CallError(FileStepEncrypting, err)
		return
	}

	if s.DeleteInput {
		s.ManagementData.CallStep(FileStepDeletingInput)
		if err = os.Remove(s.InputFilename); err != nil {
			s.ManagementData.CallError(FileStepDeletingInput, err)
			return
		}
	}

	s.ManagementData.CallStep(StepDone)
}

func (s FileData) Decrypt() {
	s.ManagementData.CallStep(FileStepCheckingExistence)
	_, err := os.Stat(s.InputFilename)
	if err != nil {
		s.ManagementData.CallError(FileStepCheckingExistence, err)
		return
	}

	if len(*s.OutputFilename) == 0 {
		*s.OutputFilename = strings.TrimSuffix(s.InputFilename, ".enc")
	}

	if *s.OutputFilename == s.InputFilename {
		s.ManagementData.CallStep(FileStepInvalidInput)
		return
	}

	password := s.PasswordReader(FileStepReadingPassword)
	if len(password) == 0 {
		return
	}

	s.ManagementData.CallStep(FileStepOpeningInput)
	inputFile, err := os.Open(s.InputFilename)
	if err != nil {
		s.ManagementData.CallError(FileStepOpeningInput, err)
		return
	}

	s.ManagementData.CallStep(FileStepCreatingOutput)
	outputFile, err := os.Create(*s.OutputFilename)
	if err != nil {
		s.ManagementData.CallError(FileStepCreatingOutput, err)
		return
	}

	s.ManagementData.CallStep(FileStepDecrypting)
	if err = cipher.StreamDecrypt(inputFile, outputFile, password); err != nil {
		os.Remove(*s.OutputFilename)
		s.ManagementData.CallError(FileStepDecrypting, err)
		return
	}

	s.ManagementData.CallStep(StepDone)
}
