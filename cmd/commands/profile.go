package commands

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/cmd/cmdutility"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/management"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"p"},
	Short:   "Profile management",
	Run:     func(cmd *cobra.Command, args []string) {},
}

var profileAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "n", "new"},
	Args:    cobra.ExactArgs(1),
	Short:   "Adds a new profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		management.ProfileData{
			ProfileName: profileName,
			PasswordReader: func(step management.ManagementStep) string {
				password, err := cmdutility.NewPasswordPrompt("New password")
				if err != nil {
					cmdutility.LogError("Failed reading password", err)
					return ""
				}
				return password
			},
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.ProfileStepProfileExists:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s already exists.", profileName)
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Profile %s has been created.", profileName)
					case management.ProfileStepDefaultChanged:
						cmdutility.LogColor(cmdutility.Green, "Default profile changed to %s", profileName)
					}
				},
			},
		}.Add()
	},
}

var profileRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "rem", "del", "delete"},
	Args:    cobra.ExactArgs(1),
	Short:   "Removes a profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		management.ProfileData{
			ProfileName: profileName,
			YesNoPrompt: func(step management.ManagementStep) bool {
				remove, err := cmdutility.YesNoPrompt(fmt.Sprintf("Are you sure to delete profile %s?", profileName))
				if err != nil {
					cmdutility.LogError("Reading input failed", err)
					return false
				}
				return remove
			},
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.ProfileStepDoesntExist:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Profile %s has been removed.", profileName)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepRemoving:
						cmdutility.LogError("Failed to remove profile", err)
					}
				},
			},
		}.Remove()
	},
}

var profileRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rn", "ren"},
	Args:    cobra.ExactArgs(2),
	Short:   "Renames a profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		newName := args[1]
		management.ProfileData{
			ProfileName:    profileName,
			NewProfileName: newName,
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.ProfileStepDoesntExist:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
					case management.ProfileStepProfileExists:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s already exists.", newName)
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Profile %s has been renamed to %s", profileName, newName)
					case management.ProfileStepDefaultChanged:
						cmdutility.LogColor(cmdutility.Green, "Default profile changed to %s", newName)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepRenaming:
						cmdutility.LogError("Profile rename failed", err)
					}
				},
			},
		}.Rename()
	},
}

var profilePasswdCmd = &cobra.Command{
	Use:     "passwd",
	Aliases: []string{"pw"},
	Args:    cobra.ExactArgs(1),
	Short:   "Changes profile password",
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		management.ProfileData{
			ProfileName: profileName,
			PasswordReader: func(step management.ManagementStep) string {
				if step == management.ProfileStepReadingPassword {
					password, err := cmdutility.PasswordPrompt("Profile password")
					if err != nil {
						cmdutility.LogError("Failed reading password", err)
						return ""
					}

					return password
				}

				newPassword, err := cmdutility.NewPasswordPrompt("New password")
				if err != nil {
					cmdutility.LogError("Failed reading password", err)
					return ""
				}

				return newPassword
			},
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.ProfileStepDoesntExist:
						cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Profile's %s password updated.", profileName)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepLoadingProfile:
						cmdutility.LogError("Failed loading profile", err)
					case management.ProfileStepSaving:
						cmdutility.LogError("Failed saving profile", err)
					}
				},
			},
		}.Passwd()
	},
}

var profileDefaultCmd = &cobra.Command{
	Use:     "default",
	Aliases: []string{"d"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Sets or gets the default profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileName := ""
		if len(args) != 0 {
			profileName = args[0]
		}
		if len(profileName) == 0 {
			if len(data.Config.DefaultProfile) == 0 {
				cmdutility.LogColor(cmdutility.Green, "No default profile is set")
			} else {
				cmdutility.LogColor(cmdutility.Green, "Default profile is %s", data.Config.DefaultProfile)
			}
		} else {
			if !data.ProfileExists(profileName) {
				cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
				return
			}

			data.Config.DefaultProfile = profileName
			cmdutility.LogColor(cmdutility.Green, "Default profile changed to %s", profileName)
		}
	},
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "Lists all profiles",
	Run: func(cmd *cobra.Command, args []string) {
		management.ProfileData{
			LogList: func(profileName string) {
				fmt.Println(profileName)
			},
			ManagementData: management.ManagementData{
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepReadingProfiles:
						cmdutility.LogError("Failed getting profiles list", err)
					}
				},
			},
		}.List()
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileRemoveCmd)
	profileCmd.AddCommand(profileRenameCmd)
	profileCmd.AddCommand(profilePasswdCmd)
	profileCmd.AddCommand(profileDefaultCmd)
	profileCmd.AddCommand(profileListCmd)
}
