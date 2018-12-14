package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists profiles or sites of a profile or fields of a site",
	Long:  `Lists profiles or sites of a profile or fields of a site`,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
