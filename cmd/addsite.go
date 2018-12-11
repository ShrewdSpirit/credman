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

		userProfile, err := config.GetUserProfile(ProfilesDir, profileName)
		if userProfile == nil && err == nil {
			fmt.Printf("No such profile '%s'\n", profileName)
			return
		}
	},
}

func init() {
	addCmd.AddCommand(addsiteCmd)
	PasswordAddFlags(addsiteCmd)
	ProfileAddFlags(addsiteCmd)
	SiteAddFieldFlags(addsiteCmd, false)
}
