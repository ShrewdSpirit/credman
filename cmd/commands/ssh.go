// +build !windows

package commands

import (
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to ssh",
	Long:  `Connects to ssh using fields specified in given site.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
