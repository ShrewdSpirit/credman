package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

var listprofilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Lists all profiles",
	Long:  `Lists all profiles`,
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			fmt.Printf("Reading profiles at %s\n", ProfilesDir)
		}
		profiles, err := ioutil.ReadDir(ProfilesDir)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Profiles:")
		for _, profile := range profiles {
			if profile.IsDir() {
				continue
			}

			name := strings.TrimSuffix(profile.Name(), ".json")
			fmt.Println(name)
		}
	},
}

func init() {
	listCmd.AddCommand(listprofilesCmd)
}
