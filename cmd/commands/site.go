package commands

import (
	"fmt"
	"strings"

	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/management"
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:     "site",
	Aliases: []string{"s"},
	Short:   "Site management",
	Run:     func(cmd *cobra.Command, args []string) {},
}

var siteAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "n", "new"},
	Args:    cobra.ExactArgs(1),
	Short:   "Adds a new site",
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]
		profile, profilePassword := cmdutility.GetProfileCommandLine()
		if profile == nil {
			return
		}

		var password string
		if !siteAddNoPassword {
			var err error
			password, err = cmdutility.ParsePasswordGenerationFlags("Site new password")
			if err != nil {
				cmdutility.LogError("Site creation failed", err)
				return
			}
		}

		management.SiteData{
			SiteName:        siteName,
			SitePassword:    password,
			SiteFieldsMap:   siteFieldsMap,
			ProfilePassword: profilePassword,
			Profile:         profile,
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.SiteStepSiteExists:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s exists.", siteName)
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Site %s has been added to profile %s", siteName, profile.Name)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepSaving:
						cmdutility.LogError("Failed saving profile", err)
					}
				},
			},
		}.Add()
	},
}

var siteRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "rem", "del", "delete"},
	Args:    cobra.ExactArgs(1),
	Short:   "Removes a site",
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]
		profile, profilePassword := cmdutility.GetProfileCommandLine()
		if profile == nil {
			return
		}

		management.SiteData{
			SiteName:        siteName,
			ProfilePassword: profilePassword,
			Profile:         profile,
			YesNoPrompt: func(step management.ManagementStep) bool {
				remove, err := cmdutility.YesNoPrompt(fmt.Sprintf("Are you sure to delete site %s?", siteName))
				if err != nil {
					cmdutility.LogError("Reading input failed", err)
					return false
				}
				return remove
			},
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.SiteStepDoesntExist:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Site %s has been removed.", siteName)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepSaving:
						cmdutility.LogError("Failed saving profile", err)
					}
				},
			},
		}.Remove()
	},
}

var siteRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rn", "ren"},
	Args:    cobra.ExactArgs(2),
	Short:   "Renames a site",
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]
		newName := args[1]
		profile, profilePassword := cmdutility.GetProfileCommandLine()
		if profile == nil {
			return
		}

		management.SiteData{
			SiteName:        siteName,
			NewSiteName:     newName,
			ProfilePassword: profilePassword,
			Profile:         profile,
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.SiteStepDoesntExist:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Site %s has been renamed to %s.", siteName, newName)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepSaving:
						cmdutility.LogError("Failed saving profile", err)
					}
				},
			},
		}.Rename()
	},
}

var siteSetCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s", "u", "edit", "e"},
	Args:    cobra.ExactArgs(1),
	Short:   "Updates site fields",
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]
		profile, profilePassword := cmdutility.GetProfileCommandLine()
		if profile == nil {
			return
		}

		var password string
		if siteSetPassword {
			var err error
			password, err = cmdutility.ParsePasswordGenerationFlags("Site new password")
			if err != nil {
				cmdutility.LogError("Site creation failed", err)
				return
			}
		}

		management.SiteData{
			SiteName:        siteName,
			SitePassword:    password,
			SiteFieldsMap:   siteFieldsMap,
			SiteFieldsList:  siteFieldsDelete,
			ProfilePassword: profilePassword,
			Profile:         profile,
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.SiteStepDoesntExist:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Site %s has been updated.", siteName)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepSaving:
						cmdutility.LogError("Failed saving profile", err)
					}
				},
			},
		}.Set()
	},
}

var siteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Lists sites using optional pattern",
	Run: func(cmd *cobra.Command, args []string) {
		pattern := ""
		if len(args) == 1 {
			pattern = args[0]
		}

		profile, profilePassword := cmdutility.GetProfileCommandLine()
		if profile == nil {
			return
		}

		management.SiteData{
			SiteName:        pattern,
			ProfilePassword: profilePassword,
			Profile:         profile,
			Match: func(name, p1, p2, p3 string) {
				if name == p1 {
					fmt.Println(name)
				} else {
					cmdutility.LogColor(cmdutility.BoldRed, p1+"%s"+p3, p2)
				}
			},
			ManagementData: management.ManagementData{
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.SiteStepRegexCompile:
						cmdutility.LogError("Failed to compile pattern", err)
					}
				},
			},
		}.List()
	},
}

var siteGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "f"},
	Args:    cobra.ExactArgs(1),
	Short:   "Gets value(s) of specified field(s) or copy the first field into clipboard",
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]
		profile, profilePassword := cmdutility.GetProfileCommandLine()
		if profile == nil {
			return
		}

		management.SiteData{
			SiteName:        siteName,
			SiteFieldsList:  siteFieldsList,
			SiteCopyField:   siteGetCopy,
			ProfilePassword: profilePassword,
			Profile:         profile,
			LogFields: func(field, value string) {
				cmdutility.LogColor(cmdutility.HiGreen, strings.Title(field)+": %s", value)
			},
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.SiteStepDoesntExist:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
					case management.SiteStepSettingClipboardPassword:
						fmt.Println("Password copied to clipboard.")
					case management.SiteStepSettingClipboard:
						cmdutility.LogColor(cmdutility.Green, "%s copied to clipboard.", siteFieldsList[0])
					case management.SiteStepInvalidField:
						cmdutility.LogColor(cmdutility.BoldRed, "No field %s was found.", siteFieldsList[0])
					case management.SiteStepListingFields:
						cmdutility.LogColor(cmdutility.Green, "Fields of site %s:", siteName)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.SiteStepSettingClipboard,
						management.SiteStepSettingClipboardPassword:
						cmdutility.LogError("Failed write to clipboard", err)
					}
				},
			},
		}.Get()
	},
}

var siteFieldsMap map[string]string
var siteFieldsList []string
var siteFieldsDelete []string
var siteGetCopy bool
var siteSetPassword bool
var siteAddNoPassword bool

func init() {
	rootCmd.AddCommand(siteCmd)
	cmdutility.FlagsAddProfileName(siteCmd)

	siteCmd.AddCommand(siteAddCmd)
	siteAddCmd.Flags().BoolVarP(&siteAddNoPassword, "no-password", "n", false, "Doesn't prompt for site password. Useful for sites that you don't want any password for.")
	siteFlagsFields(siteAddCmd, false)
	cmdutility.FlagsAddPasswordOptions(siteAddCmd)

	siteCmd.AddCommand(siteRemoveCmd)
	siteCmd.AddCommand(siteRenameCmd)
	siteCmd.AddCommand(siteSetCmd)
	siteSetCmd.Flags().BoolVarP(&siteSetPassword, "password", "w", false, "Change password. Can be used with password generator or it will prompt user")
	siteSetCmd.Flags().StringSliceVarP(&siteFieldsDelete, "delete", "d", []string{}, "Deletes specified fields")
	siteFlagsFields(siteSetCmd, false)
	cmdutility.FlagsAddPasswordOptions(siteSetCmd)

	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteGetCmd)
	siteFlagsFields(siteGetCmd, true)
	siteGetCmd.Flags().BoolVarP(&siteGetCopy, "copy", "c", false, "Copy first selected field into clipboard")
}

func siteFlagsFields(cmd *cobra.Command, get bool) {
	if get {
		cmd.Flags().StringSliceVarP(&siteFieldsList, "fields", "f", make([]string, 0), "-f=Key1,Key2 ...")
	} else {
		cmd.Flags().StringToStringVarP(&siteFieldsMap, "field", "f", make(map[string]string), "-f=Key=Value ...")
	}
}
