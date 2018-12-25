package management

import (
	"os"

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

	// check output file if len == 0 output file = input without .enc if it had, or .dec if didnt had
	// read password
	// dec
	// delete original
	s.ManagementData.CallStep(StepDone)
}
