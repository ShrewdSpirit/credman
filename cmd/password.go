package cmd

import (
	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Password options help",
	Long: `Password options provide a set of settings for password generation.
Any password (either generated or given by user shouldn't be more than 32 characters)
If you set the --pgen flag on any of the bellow sub commands, a password will be generated instead of prompting one from user:
- add site
- passwd site

Using --plen option you can set the generated password's length
Using --pcase option you can specify the password's letter case mix
Using --pmix option you can specify which types of characters should be included in password`,
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(passwordCmd)
}
