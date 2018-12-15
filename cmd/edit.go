package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit site fields",
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

		if Verbose {
			fmt.Println("Setting fields")
		}

		email, _ := cmd.Flags().GetString("email")
		username, _ := cmd.Flags().GetString("username")

		fields, _ := cmd.Flags().GetStringToString("fields")

		if len(email) > 0 {
			fields["email"] = email
		}
		if len(username) > 0 {
			fields["username"] = username
		}

		setFieldsBuffer := bytes.Buffer{}
		for field, value := range fields {
			site.Fields[field] = value
			setFieldsBuffer.WriteString(field + ", ")
		}
		setFields := setFieldsBuffer.String()
		setFields = strings.TrimSuffix(setFields, ", ")

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

		fmt.Printf("Fields %s has been updated\n", setFields)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	AddFlagsProfileName(editCmd)
	AddFlagsSiteFields(editCmd, false)
}
