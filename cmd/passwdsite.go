package cmd

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/spf13/cobra"
)

var passwdsiteCmd = &cobra.Command{
	Use:   "site",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]

		profileName, _ := cmd.Flags().GetString("profile")
		if profileName == "" {
			profileName = config.AppConfig.DefaultProfile
		}
		if profileName == "" {
			fmt.Println("Create a profile first")
			return
		}

		if Verbose {
			fmt.Println("Reading profile")
		}
		profile, err := config.GetProfile(ProfilesDir, profileName)
		if profile == nil && err != nil {
			fmt.Printf("No such profile '%s'\n", profileName)
			return
		}

		profilePassword, err := PromptPassword("Profile password")
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
			fmt.Println("Checking site's existence")
		}
		site := profile.Site(siteName)
		if site == nil {
			fmt.Printf("Site '%s' doesn't exist\n", siteName)
			return
		}

		password, err := GetPassword(cmd, "New password")
		if err != nil {
			fmt.Println(err)
			return
		}

		site.Fields["password"] = password

		if Verbose {
			fmt.Println("Encrypting profile")
		}
		if err := profile.Encrypt(profilePassword); err != nil {
			fmt.Println(err)
			return
		}

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
	passwdCmd.AddCommand(passwdsiteCmd)
	AddFlagsPassword(passwdsiteCmd)
	AddFlagsProfileName(passwdsiteCmd)
}
