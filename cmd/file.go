package cmd

import (
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File encryption",
}

var fileEncryptCmd = &cobra.Command{
	Use:     "encrypt",
	Aliases: []string{"enc", "e"},
	Args:    cobra.RangeArgs(1, 2),
	Short:   "Encrypt file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			fileEncrypt(args[0], "")
		} else {
			fileEncrypt(args[0], args[1])
		}
	},
}

var fileDecryptCmd = &cobra.Command{
	Use:     "decrypt",
	Aliases: []string{"dec", "d"},
	Args:    cobra.ExactArgs(1),
	Short:   "Decrypt file",
	Run: func(cmd *cobra.Command, args []string) {
		fileDecrypt(args[0])
	},
}

var fileListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "List profile files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fileList("")
		} else {
			fileList(args[0])
		}
	},
}

var fileDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d", "remove", "rm"},
	Args:    cobra.ExactArgs(1),
	Short:   "Delete profile file",
	Run: func(cmd *cobra.Command, args []string) {
		fileDelete(args[0])
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
	fileCmd.AddCommand(fileEncryptCmd)
	fileCmd.AddCommand(fileDecryptCmd)
	fileCmd.AddCommand(fileListCmd)
	fileCmd.AddCommand(fileDeleteCmd)
}

func fileEncrypt(filename, name string) {}

func fileDecrypt(name string) {}

func fileList(pattern string) {}

func fileDelete(name string) {}
