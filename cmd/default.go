package cmd

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Sets or gets the default profile",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("Default profile is '%s'\n", config.AppConfig.DefaultProfile)
		} else {
			profileName := args[0]

			if Verbose {
				fmt.Println("Checking profile's existence")
			}
			_, err := config.GetProfile(ProfilesDir, profileName)
			if err != nil {
				fmt.Printf("Failed reading profile '%s': %s\n", profileName, err)
				return
			}

			config.AppConfig.DefaultProfile = profileName
			fmt.Printf("Default profile changed to '%s'\n", config.AppConfig.DefaultProfile)
		}
	},
}

func init() {
	rootCmd.AddCommand(defaultCmd)
}
