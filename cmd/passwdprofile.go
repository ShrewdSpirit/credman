package cmd

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/spf13/cobra"
)

var passwdprofileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Changes profile password",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]

		if Verbose {
			fmt.Println("Reading profile")
		}
		profile, err := config.GetProfile(ProfilesDir, profileName)
		if profile == nil && err != nil {
			fmt.Printf("No such profile '%s'\n", profileName)
			return
		}

		profilePassword, err := PromptPassword("Current password")
		if err != nil {
			fmt.Println(err)
			return
		}

		if !profile.CheckPassword(profilePassword) {
			fmt.Println("Wrong password")
			return
		}

		newPassword, err := GetPassword(cmd, "New password")
		if err != nil {
			fmt.Println(err)
			return
		}

		if Verbose {
			fmt.Println("Decrypting profile")
		}
		err = profile.Decrypt(profilePassword)
		if err != nil {
			fmt.Println(err)
			return
		}

		if Verbose {
			fmt.Println("Encrypting profile")
		}
		if err := profile.Encrypt(newPassword); err != nil {
			fmt.Println(err)
			return
		}

		profile.Hash = config.Hash(newPassword)

		if Verbose {
			fmt.Println("Saving profile")
		}
		if err := profile.Save(ProfilesDir); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Password changed")
	},
}

func init() {
	passwdCmd.AddCommand(passwdprofileCmd)
}
