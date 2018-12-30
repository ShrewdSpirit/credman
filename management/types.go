package management

type ManagementStep int

const (
	StepDone ManagementStep = iota
	StepFailed

	ProfileStepCheckingExistence
	ProfileStepProfileExists
	ProfileStepReadingPassword
	ProfileStepCreating
	ProfileStepSaving
	ProfileStepDefaultChanged
	ProfileStepDoesntExist
	ProfileStepRemovePrompt
	ProfileStepRemoving
	ProfileStepRenaming
	ProfileStepNewPassword
	ProfileStepLoadingProfile
	ProfileStepReadingProfiles

	SiteStepAdding
	SiteStepSiteExists
	SiteStepCreating
	SiteStepSettingFields
	SiteStepCheckingExistence
	SiteStepDoesntExist
	SiteStepRemoving
	SiteStepRenaming
	SiteStepUpdatingFields
	SiteStepDeletingFields
	SiteStepRegexCompile
	SiteStepSettingClipboardPassword
	SiteStepSettingClipboard
	SiteStepInvalidField
	SiteStepListingFields
	SiteStepAddingTags
	SiteStepRemovingTags

	FileStepCheckingExistence
	FileStepInputDoesntExist
	FileStepInvalidOutput
	FileStepEncrypting
	FileStepDecrypting
	FileStepReadingPassword
	FileStepOpeningInput
	FileStepCreatingOutput
	FileStepDeletingInput
	FileStepInvalidInput

	RestoreStepReadingPassword
	RestoreStepEncrypting
	RestoreStepDecrypting
	RestoreStepClipboardPassword
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
