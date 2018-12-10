package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "credman",
	Short: "Simple yet powerful credential/password manager",
	Long: `Simple yet powerful credential/password manager with remote sync.

Use 'credman help password' to see how to use password options.`,
	Version: "0.1.0",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Prints verbose information")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
