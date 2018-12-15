package cmd

import (
	"github.com/spf13/cobra"
)

var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Sets profile/site password",
}

func init() {
	rootCmd.AddCommand(passwdCmd)
}
