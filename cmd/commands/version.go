package commands

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Short:   "Prints version information",
	Use:     "version",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", data.Version)
		fmt.Printf("Git commit hash: %s\n", data.GitCommit)
	},
}

func init() {
	rootCmd.AddCommand(cmdVersion)
}
