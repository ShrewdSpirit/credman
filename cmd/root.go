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

Use 'credman help password' to see how to use password options.
Use 'credman help fields to see what fields you can set for a site.`,
	Version: "0.2.0",
}

func init() {
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
