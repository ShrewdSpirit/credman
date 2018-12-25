package management

type ManagementStep int

const (
	StepDone   = 0
	StepFailed = 1

	ProfileStepCheckingExistence ManagementStep = 2
	ProfileStepProfileExists     ManagementStep = 3
	ProfileStepReadingPassword   ManagementStep = 4
	ProfileStepCreating          ManagementStep = 5
	ProfileStepSaving            ManagementStep = 6
	ProfileStepDefaultChanged    ManagementStep = 8
	ProfileStepDoesntExist       ManagementStep = 9
	ProfileStepRemovePrompt      ManagementStep = 10
	ProfileStepRemoving          ManagementStep = 11
	ProfileStepRenaming          ManagementStep = 12
	ProfileStepNewPassword       ManagementStep = 13
	ProfileStepLoadingProfile    ManagementStep = 14
	ProfileStepReadingProfiles   ManagementStep = 15

	SiteStepAdding                   ManagementStep = 16
	SiteStepSiteExists               ManagementStep = 17
	SiteStepCreating                 ManagementStep = 18
	SiteStepSettingFields            ManagementStep = 19
	SiteStepCheckingExistence        ManagementStep = 20
	SiteStepDoesntExist              ManagementStep = 21
	SiteStepRemoving                 ManagementStep = 22
	SiteStepRenaming                 ManagementStep = 23
	SiteStepUpdatingFields           ManagementStep = 24
	SiteStepDeletingFields           ManagementStep = 25
	SiteStepRegexCompile             ManagementStep = 26
	SiteStepSettingClipboardPassword ManagementStep = 27
	SiteStepSettingClipboard         ManagementStep = 28
	SiteStepInvalidField             ManagementStep = 29
	SiteStepListingFields            ManagementStep = 30

	FileStepCheckingExistence ManagementStep = 31
	FileStepInputDoesntExist  ManagementStep = 32
	FileStepInvalidOutput     ManagementStep = 33
	FileStepEncrypting        ManagementStep = 34
	FileStepDecrypting        ManagementStep = 35
	FileStepReadingPassword   ManagementStep = 36
	FileStepOpeningInput      ManagementStep = 37
	FileStepCreatingOutput    ManagementStep = 38
	FileStepDeletingInput     ManagementStep = 39
)

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
type LogString2Func func(string, string)
type LogMatch func(name, p1, p2, p3 string)
