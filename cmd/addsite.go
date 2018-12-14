package cmd

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/spf13/cobra"
)

var addsiteCmd = &cobra.Command{
	Use:     "site",
	Short:   "Adds a new site",
	Long:    "Adds a new site",
	Example: `credman add site website.com --email="mymail@m.com"`,
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

		fmt.Printf("Adding site '%s' to profile '%s'\n", siteName, profileName)

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
		if profile.SiteExist(siteName) {
			fmt.Printf("Site '%s' exists\n", siteName)
			return
		}

		password, err := GetPassword(cmd, "Site password")
		if err != nil {
			fmt.Println(err)
			return
		}

		if Verbose {
			fmt.Println("Creating site")
		}
		site := config.NewSite(siteName, password, cmd)

		if Verbose {
			fmt.Println("Adding site to profile's list")
		}
		profile.AddSite(site)

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

		fmt.Printf("Site '%s' has been added\n", siteName)
	},
}

func init() {
	addCmd.AddCommand(addsiteCmd)
	PasswordAddFlags(addsiteCmd)
	ProfileAddFlags(addsiteCmd)
	SiteAddFieldFlags(addsiteCmd, false)
}
