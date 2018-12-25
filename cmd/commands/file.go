package commands

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

var fileOutput string
var fileDeleteOriginal bool
var fileNoProfile bool

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.AddCommand(fileEncryptCmd)
	fileFlagsAddOutput(fileEncryptCmd)
	fileFlagsAddNoProfile(fileEncryptCmd)
	fileEncryptCmd.Flags().BoolVarP(&fileDeleteOriginal, "delete-original", "d", false, "Deletes original file after encryption")
	FlagsAddPasswordOptions(fileEncryptCmd)
	FlagsAddProfileName(fileEncryptCmd)

	fileCmd.AddCommand(fileDecryptCmd)
	fileFlagsAddOutput(fileDecryptCmd)
	fileFlagsAddNoProfile(fileDecryptCmd)
	FlagsAddProfileName(fileDecryptCmd)

	fileCmd.AddCommand(fileListCmd)
	FlagsAddProfileName(fileListCmd)

	fileCmd.AddCommand(fileDeleteCmd)
	FlagsAddProfileName(fileDeleteCmd)
}

func fileFlagsAddOutput(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&fileOutput, "output", "o", "", "Output file for encryption/decryption")
}

func fileFlagsAddNoProfile(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&fileNoProfile, "no-profile", "n", false, "Doesn't use a profile for encryption/decryption")
}

func fileEncrypt(filename, name string) {
	// if fileNoProfile {
	// 	_, err := os.Stat(filename)
	// 	if err != nil {
	// 		utility.LogError("Invalid source file", err)
	// 		return
	// 	}

	// 	password, err := ParsePasswordGenerationFlags("Encryption password")
	// 	if err != nil {
	// 		utility.LogError("Password reading failed", err)
	// 		return
	// 	}

	// 	utility.LogColor(utility.Green, "Reading file %s", filename)
	// 	data, err := ioutil.ReadFile(filename)
	// 	if err != nil {
	// 		utility.LogError("Reading file failed", err)
	// 		return
	// 	}

	// 	utility.LogColor(utility.Green, "Encrypting")
	// 	encrypted, err := utility.Encrypt([]byte(password), data)
	// 	if err != nil {
	// 		utility.LogError("Encryption failed", err)
	// 	}

	// 	utility.LogColor(utility.Green, "Writing to %s", name)
	// 	if err := ioutil.WriteFile(name, encrypted, os.ModePerm); err != nil {
	// 		utility.LogError("Writing file failed", err)
	// 		return
	// 	}
	// } else {
	// 	fmt.Println("Not implemented")
	// }

	// if fileDeleteOriginal {
	// 	utility.LogColor(utility.Green, "Deleting original file %s", filename)
	// 	if err := os.Remove(filename); err != nil {
	// 		utility.LogError("Removing original file failed", err)
	// 		return
	// 	}
	// }
}

func fileDecrypt(name string) {

}

func fileList(pattern string) {
}

func fileDelete(name string) {
}
