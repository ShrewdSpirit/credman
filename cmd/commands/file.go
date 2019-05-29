package commands

import (
	"os"
	"strings"

	"github.com/ShrewdSpirit/credman/cipher"
	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:     "file",
	Aliases: []string{"f"},
	Short:   "File encryption",
}

var fileEncryptCmd = &cobra.Command{
	Use:     "encrypt",
	Aliases: []string{"enc", "e"},
	Args:    cobra.ExactArgs(1),
	Short:   "Encrypt file",
	Run:     fileEncrypt,
}

var fileDecryptCmd = &cobra.Command{
	Use:     "decrypt",
	Aliases: []string{"dec", "d"},
	Args:    cobra.ExactArgs(1),
	Short:   "Decrypt file",
	Run:     fileDecrypt,
}

var outputFilename string
var fileDeleteOriginal bool
var fileNoNewSite bool
var fileStore bool

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.AddCommand(fileEncryptCmd)
	fileFlagsAddOutput(fileEncryptCmd)
	fileEncryptCmd.Flags().BoolVarP(&fileDeleteOriginal, "delete-original", "d", false, "Deletes original file after encryption")
	fileEncryptCmd.Flags().BoolVarP(&fileNoNewSite, "no-site", "n", false, "Doesn't store the file information as a site")
	fileEncryptCmd.Flags().BoolVarP(&fileStore, "store", "s", false, "Store file's encrypted content in file's site")

	fileCmd.AddCommand(fileDecryptCmd)
	fileFlagsAddOutput(fileDecryptCmd)
}

func fileFlagsAddOutput(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outputFilename, "output", "o", "", "Output file for encryption/decryption")
}

func fileEncrypt(cmd *cobra.Command, args []string) {
	inputFilename := args[0]

	if _, err := os.Stat(inputFilename); err != nil {
		cmdutility.LogError("File doesn't exist", err)
		return
	}

	if len(outputFilename) == 0 {
		outputFilename = inputFilename + ".enc"
	}

	password, err := cmdutility.NewPasswordPrompt("New password")
	if err != nil {
		cmdutility.LogError("Failed reading password", err)
		return
	}

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		cmdutility.LogError("Failed to open input file", err)
		return
	}

	if !fileStore {
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			cmdutility.LogError("Failed to create output file", err)
			return
		}

		cmdutility.LogColor(cmdutility.Green, "Encrypting file %s", inputFilename)
		if err := cipher.StreamEncrypt(inputFile, outputFile, password); err != nil {
			os.Remove(outputFilename)
			cmdutility.LogError("Failed to encrypt file", err)
			return
		}
	} // TODO: store in site

	if fileDeleteOriginal {
		cmdutility.LogColor(cmdutility.Green, "Deleting original file %s", inputFilename)
		if err := os.Remove(inputFilename); err != nil {
			cmdutility.LogError("Failed to delete original file", err)
			return
		}
	}

	cmdutility.LogColor(cmdutility.Green, "Encryption done. Wrote data to %s", outputFilename)
}

func fileDecrypt(cmd *cobra.Command, args []string) {
	inputFilename := args[0]

	if _, err := os.Stat(inputFilename); err != nil {
		cmdutility.LogError("File doesn't exist", err)
		return
	}

	if len(outputFilename) == 0 {
		outputFilename = strings.TrimSuffix(inputFilename, ".enc")
	}

	if outputFilename == inputFilename {
		cmdutility.LogColor(cmdutility.BoldRed, "Input file name %s is same as output %s", inputFilename, outputFilename)
		return
	}

	password, err := cmdutility.PasswordPrompt("File password")
	if err != nil {
		cmdutility.LogError("Failed reading password", err)
		return
	}

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		cmdutility.LogError("Failed to open input file", err)
		return
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		cmdutility.LogError("Failed to create output file", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Decrypting file %s", inputFilename)
	if err := cipher.StreamDecrypt(inputFile, outputFile, password); err != nil {
		os.Remove(outputFilename)
		cmdutility.LogError("Failed to decrypt file", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Decryption done. Wrote data to %s", outputFilename)
}
