package cmd

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/spf13/cobra"
)

var addprofileCmd = &cobra.Command{
	Use:     "profile",
	Short:   "Adds a new profile",
	Long:    "Adds a new profile",
	Example: "credman add profile myprofile",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		_, err := config.GetUserProfile(ProfilesDir, profileName)
		if err == nil {
			fmt.Printf("Profile %s exists\n", profileName)
			return
		}

		password, err := GetPassword(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		config.NewUserProfile(ProfilesDir, profileName, password)
	},
}

func init() {
	addCmd.AddCommand(addprofileCmd)
	PasswordAddFlags(addprofileCmd)
	ProfileAddFlags(addprofileCmd)
}
