package commands

import (
	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/management"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File encryption",
}

var fileEncryptCmd = &cobra.Command{
	Use:     "encrypt",
	Aliases: []string{"enc", "e"},
	Args:    cobra.ExactArgs(1),
	Short:   "Encrypt file",
	Run: func(cmd *cobra.Command, args []string) {
		inputFilename := args[0]

		management.FileData{
			InputFilename:  inputFilename,
			OutputFilename: fileOutput,
			PasswordReader: func(step management.ManagementStep) string {
				password, err := cmdutility.NewPasswordPrompt("New password")
				if err != nil {
					cmdutility.LogError("Failed reading password", err)
					return ""
				}
				return password
			},
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {

				},
				OnError: func(step management.ManagementStep, err error) {

				},
			},
		}.Encrypt()
	},
}

var fileDecryptCmd = &cobra.Command{
	Use:     "decrypt",
	Aliases: []string{"dec", "d"},
	Args:    cobra.ExactArgs(1),
	Short:   "Decrypt file",
	Run: func(cmd *cobra.Command, args []string) {
		inputFilename := args[0]

		management.FileData{
			InputFilename:  inputFilename,
			OutputFilename: fileOutput,
			PasswordReader: func(step management.ManagementStep) string {
				password, err := cmdutility.PasswordPrompt("File password")
				if err != nil {
					cmdutility.LogError("Failed reading password", err)
					return ""
				}
				return password
			},
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {

				},
				OnError: func(step management.ManagementStep, err error) {

				},
			},
		}.Decrypt()
	},
}

var fileOutput string
var fileDeleteOriginal bool

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.AddCommand(fileEncryptCmd)
	fileFlagsAddOutput(fileEncryptCmd)
	fileEncryptCmd.Flags().BoolVarP(&fileDeleteOriginal, "delete-original", "d", false, "Deletes original file after encryption")

	fileCmd.AddCommand(fileDecryptCmd)
	fileFlagsAddOutput(fileDecryptCmd)
}

func fileFlagsAddOutput(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&fileOutput, "output", "o", "", "Output file for encryption/decryption")
}
