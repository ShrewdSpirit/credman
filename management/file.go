package management

type FileData struct {
	ManagementData
	InputFilename  string
	OutputFilename string
	PasswordReader GetStringFunc
}

func (s FileData) Encrypt() {
	s.ManagementData.CallStep(StepDone)
}

func (s FileData) Decrypt() {
	s.ManagementData.CallStep(StepDone)
}
