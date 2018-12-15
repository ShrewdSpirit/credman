package cmd

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/spf13/cobra"
)

var listsitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Lists sites of default profile or given profile's name",
	Run: func(cmd *cobra.Command, args []string) {
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
		if err != nil {
			fmt.Printf("Failed reading profile '%s': %s\n", profileName, err)
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

		fmt.Printf("Sites of '%s':\n", profileName)
		sorted := profile.Sites.Sort()
		for _, site := range sorted {
			fmt.Println(site.Name)
		}
	},
}

func init() {
	listCmd.AddCommand(listsitesCmd)
	AddFlagsProfileName(listsitesCmd)
}
