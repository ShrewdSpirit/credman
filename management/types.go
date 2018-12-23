package management

type ManagementStep int

type ManagementData struct {
	OnError LogErrorFunc
	OnStep  LogStepFunc
}

func (s ManagementData) CallStep(step ManagementStep) {
	if s.OnStep != nil {
		s.OnStep(step)
	}
}

func (s ManagementData) CallError(step ManagementStep, err error) {
	if s.OnError != nil {
		s.OnError(step, err)
	}
}

type GetStringFunc func(ManagementStep) string
type GetYesNoPromptFunc func(ManagementStep) bool
type LogErrorFunc func(ManagementStep, error)
type LogStepFunc func(ManagementStep)
type LogStringFunc func(string)
