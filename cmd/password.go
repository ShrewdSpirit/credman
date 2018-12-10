package cmd

import (
	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:    "password",
	Short:  "A brief description of your command",
	Long:   `Password options help`,
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(passwordCmd)
}
