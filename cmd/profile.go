package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/utility"
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
		profileAdd(args[0])
	},
}

var profileRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "rem", "del", "delete"},
	Args:    cobra.ExactArgs(1),
	Short:   "Removes a profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileRemove(args[0])
	},
}

var profileRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rn", "ren"},
	Args:    cobra.ExactArgs(2),
	Short:   "Renames a profile",
	Run: func(cmd *cobra.Command, args []string) {
		profileRename(args[0], args[1])
	},
}

var profilePasswdCmd = &cobra.Command{
	Use:     "passwd",
	Aliases: []string{"pw"},
	Args:    cobra.ExactArgs(1),
	Short:   "Changes profile password",
	Run: func(cmd *cobra.Command, args []string) {
		profilePasswd(args[0])
	},
}

var profileDefaultCmd = &cobra.Command{
	Use:     "default",
	Aliases: []string{"d"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Sets or gets the default profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			profileDefault("")
		} else {
			profileDefault(args[0])
		}
	},
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "Lists all profiles",
	Run: func(cmd *cobra.Command, args []string) {
		profileList()
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

func profileAdd(profileName string) {
	if data.ProfileExists(profileName) {
		utility.LogColor(utility.BoldHiYellow, "Profile %s already exists.", profileName)
		return
	}

	password, err := utility.NewPasswordPrompt("New password")
	if err != nil {
		utility.LogError("Failed reading password", err)
		return
	}

	profile := data.NewProfile(profileName, password)
	profile.ProfileMeta.CreationDate = time.Now().UnixNano()

	if err = os.Mkdir(path.Join(data.ProfilesDir, profileName), os.ModePerm); err != nil {
		utility.LogError("Failed creating profile directory", err)
		return
	}

	if err = profile.Save(password); err != nil {
		utility.LogError("Failed saving profile", err)
		return
	}

	utility.LogColor(utility.Green, "Profile %s has been created.", profileName)

	if data.Config.DefaultProfile == "" {
		data.Config.DefaultProfile = profileName
		utility.LogColor(utility.Green, "Default profile changed to %s", profileName)
	}
}

func profileRemove(profileName string) {
	if !data.ProfileExists(profileName) {
		utility.LogColor(utility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
		return
	}

	remove, err := utility.YesNoPrompt(fmt.Sprintf("Are you sure to delete profile %s?", profileName))
	if err != nil {
		utility.LogError("Reading input failed", err)
		return
	}

	if !remove {
		return
	}

	if err = os.RemoveAll(path.Join(data.ProfilesDir, profileName)); err != nil {
		utility.LogError("Failed to remove profile", err)
		return
	}

	data.Config.DefaultProfile = ""
	utility.LogColor(utility.Green, "Profile %s has been removed.", profileName)
}

func profileRename(profileName, newName string) {
	if !data.ProfileExists(profileName) {
		utility.LogColor(utility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
		return
	}

	if data.ProfileExists(newName) {
		utility.LogColor(utility.BoldHiYellow, "Profile %s already exists.", newName)
		return
	}

	ppath := path.Join(data.ProfilesDir, profileName)
	npath := path.Join(data.ProfilesDir, newName)
	if err := os.Rename(ppath, npath); err != nil {
		utility.LogError("Profile rename failed", err)
		return
	}

	utility.LogColor(utility.Green, "Profile %s has been renamed to %s", profileName, newName)
	if data.Config.DefaultProfile == profileName {
		data.Config.DefaultProfile = newName
		utility.LogColor(utility.Green, "Default profile changed to %s", newName)
	}
}

func profilePasswd(profileName string) {
	if !data.ProfileExists(profileName) {
		utility.LogColor(utility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
		return
	}

	password, err := utility.PasswordPrompt("Profile password")
	if err != nil {
		utility.LogError("Failed reading password", err)
		return
	}

	profile, err := data.LoadProfile(profileName, password)
	if err != nil {
		utility.LogError("Failed loading profile", err)
		return
	}

	newPassword, err := utility.NewPasswordPrompt("New password")
	if err != nil {
		utility.LogError("Failed reading password", err)
		return
	}

	profile.Hash = utility.Hash(newPassword)
	if err := profile.Save(newPassword); err != nil {
		utility.LogError("Failed saving profile", err)
		return
	}

	utility.LogColor(utility.Green, "Profile's %s password updated.", profileName)
}

func profileDefault(profileName string) {
	if len(profileName) == 0 {
		if len(data.Config.DefaultProfile) == 0 {
			utility.LogColor(utility.Green, "No default profile is set")
		} else {
			utility.LogColor(utility.Green, "Default profile is %s", data.Config.DefaultProfile)
		}
	} else {
		if !data.ProfileExists(profileName) {
			utility.LogColor(utility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
			return
		}

		data.Config.DefaultProfile = profileName
		utility.LogColor(utility.Green, "Default profile changed to %s", profileName)
	}
}

func profileList() {
	profiles, err := ioutil.ReadDir(data.ProfilesDir)
	if err != nil {
		utility.LogError("Failed getting profiles list", err)
		return
	}

	for _, profile := range profiles {
		if profile.IsDir() {
			fmt.Println(profile.Name())
		}
	}
}
