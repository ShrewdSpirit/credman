package cmd

import (
	"github.com/spf13/cobra"
)

var addsiteCmd = &cobra.Command{
	Use:     "site",
	Short:   "Adds a new site",
	Long:    "Adds a new site",
	Example: `credman add site website.com --email="mymail@m.com"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO:
	},
}

func init() {
	addCmd.AddCommand(addsiteCmd)
	PasswordAddFlags(addsiteCmd)
	ProfileAddFlags(addsiteCmd)
}
