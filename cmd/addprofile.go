package cmd

import (
	"github.com/spf13/cobra"
)

var addprofileCmd = &cobra.Command{
	Use:     "profile",
	Short:   "Adds a new profile",
	Long:    "Adds a new profile",
	Example: "credman add profile myprofile",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO:
	},
}

func init() {
	addCmd.AddCommand(addprofileCmd)
	PasswordAddFlags(addprofileCmd)
	ProfileAddFlags(addprofileCmd)
}
