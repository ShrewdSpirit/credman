package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/ShrewdSpirit/credman/cmd/cmdutility"

	"github.com/ShrewdSpirit/credman/data"
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
	Run:     profileAdd,
}

var profileRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "rem", "del", "delete"},
	Args:    cobra.ExactArgs(1),
	Short:   "Removes a profile",
	Run:     profileRemove,
}

var profileRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rn", "ren"},
	Args:    cobra.ExactArgs(2),
	Short:   "Renames a profile",
	Run:     profileRename,
}

var profilePasswdCmd = &cobra.Command{
	Use:     "passwd",
	Aliases: []string{"pw"},
	Args:    cobra.ExactArgs(1),
	Short:   "Changes profile password",
	Run:     profilePasswd,
}

var profileDefaultCmd = &cobra.Command{
	Use:     "default",
	Aliases: []string{"d"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Sets or gets the default profile",
	Run:     profileDefault,
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "Lists all profiles",
	Run:     profileList,
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

func profileAdd(cmd *cobra.Command, args []string) {
	profileName := args[0]

	if data.ProfileExists(profileName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s already exists.", profileName)
		return
	}

	password, err := cmdutility.NewPasswordPrompt("New password")
	if err != nil {
		cmdutility.LogError("Failed reading password", err)
		return
	}

	profile := data.NewProfile(profileName)
	if err := profile.Save(password); err != nil {
		cmdutility.LogError("Failed saving profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Profile %s has been created.", profileName)

	if data.Config.DefaultProfile == "" {
		data.Config.DefaultProfile = profileName
		cmdutility.LogColor(cmdutility.Green, "Default profile changed to %s", profileName)
	}
}

func profileDefault(cmd *cobra.Command, args []string) {
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
}

func profileList(cmd *cobra.Command, args []string) {
	profiles, err := data.ListProfiles()
	if err != nil {
		cmdutility.LogError("Failed getting profiles list", err)
		return
	}

	for _, profile := range profiles {
		fmt.Println(profile)
	}
}

func profilePasswd(cmd *cobra.Command, args []string) {
	profileName := args[0]

	if !data.ProfileExists(profileName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
		return
	}

	password, err := cmdutility.PasswordPrompt("Profile password")
	if err != nil {
		cmdutility.LogError("Failed reading password", err)
		return
	}

	profile, err := data.LoadProfile(profileName, password)
	if err != nil {
		cmdutility.LogError("Failed loading profile", err)
		return
	}

	newPassword, err := cmdutility.NewPasswordPrompt("New password")
	if err != nil {
		cmdutility.LogError("Failed reading password", err)
		return
	}

	if err := profile.Save(newPassword); err != nil {
		cmdutility.LogError("Failed saving profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Profile's %s password updated.", profileName)
}

func profileRemove(cmd *cobra.Command, args []string) {
	profileName := args[0]

	if !data.ProfileExists(profileName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
		return
	}

	remove, err := cmdutility.YesNoPrompt(fmt.Sprintf("Are you sure to delete profile %s?", profileName))
	if err != nil {
		cmdutility.LogError("Reading input failed", err)
		return
	}

	if !remove {
		return
	}

	if err := os.RemoveAll(path.Join(data.ProfilesDir, profileName)); err != nil {
		cmdutility.LogError("Failed to remove profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Profile %s has been removed.", profileName)

	if data.Config.DefaultProfile == profileName {
		data.Config.DefaultProfile = ""
	}
}

func profileRename(cmd *cobra.Command, args []string) {
	profileName := args[0]
	newName := args[1]

	if !data.ProfileExists(profileName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
		return
	}

	if data.ProfileExists(newName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Profile %s already exists.", newName)
		return
	}

	ppath := path.Join(data.ProfilesDir, profileName)
	npath := path.Join(data.ProfilesDir, newName)
	if err := os.Rename(ppath, npath); err != nil {
		cmdutility.LogError("Profile rename failed", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Profile %s has been renamed to %s", profileName, newName)

	if data.Config.DefaultProfile == profileName {
		data.Config.DefaultProfile = newName
		cmdutility.LogColor(cmdutility.Green, "Default profile changed to %s", newName)
	}
}
