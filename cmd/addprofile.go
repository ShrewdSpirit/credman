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

		if Verbose {
			fmt.Println("Checking profile's existence")
		}
		profile, err := config.GetProfile(ProfilesDir, profileName)
		if profile != nil && err == nil {
			fmt.Printf("Profile '%s' exists\n", profileName)
			return
		} else if err != nil {
			fmt.Printf("Failed checking profile '%s': %s\n", profileName, err)
			return
		}

		password, err := GetPassword(cmd, "")
		if err != nil {
			fmt.Println(err)
			return
		}

		if Verbose {
			fmt.Println("Creating new profile")
		}
		config.NewProfile(ProfilesDir, profileName, password)
		fmt.Printf("Profile '%s' has been added\n", profileName)

		if len(config.AppConfig.DefaultProfile) == 0 {
			config.AppConfig.DefaultProfile = profileName
			fmt.Printf("Profile '%s' has been set as default profile\n", profileName)
		}
	},
}

func init() {
	addCmd.AddCommand(addprofileCmd)
	ProfileAddFlags(addprofileCmd)
}
