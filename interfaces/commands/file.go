package commands

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ShrewdSpirit/credman/cipher"
	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/interfaces/commands/cmdutility"
	"github.com/ShrewdSpirit/credman/utils"
	"github.com/google/uuid"
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

var fileDeleteOriginal bool
var fileStore bool

func init() {
	rootCmd.AddCommand(fileCmd)
	cmdutility.FlagsAddProfileName(fileCmd)

	fileCmd.AddCommand(fileEncryptCmd)
	cmdutility.FlagsAddOutput(fileEncryptCmd, "Output file for encryption/decryption")
	fileEncryptCmd.Flags().BoolVarP(&fileDeleteOriginal, "delete-original", "d", false, "Deletes original file after encryption")
	fileEncryptCmd.Flags().BoolVarP(&fileStore, "store", "s", false, "Store file's encrypted content in file's site")

	fileCmd.AddCommand(fileDecryptCmd)
	cmdutility.FlagsAddOutput(fileDecryptCmd, "Output file for encryption/decryption")
}

func fileEncrypt(cmd *cobra.Command, args []string) {
	inputFilename := args[0]

	stat, err := os.Stat(inputFilename)
	if err != nil {
		cmdutility.LogError("File doesn't exist", err)
		return
	}

	if len(cmdutility.FlagOutput) == 0 {
		cmdutility.FlagOutput = inputFilename + ".enc"
	}

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		cmdutility.LogError("Failed to open input file", err)
		return
	}

	if !fileStore {
		password, err := cmdutility.NewPasswordPrompt("File password")
		if err != nil {
			cmdutility.LogError("Failed reading password", err)
			return
		}

		outputFile, err := os.Create(cmdutility.FlagOutput)
		if err != nil {
			cmdutility.LogError("Failed to create output file", err)
			return
		}

		cmdutility.LogColor(cmdutility.Green, "Encrypting file %s", inputFilename)
		if err := cipher.StreamEncrypt(inputFile, outputFile, password); err != nil {
			os.Remove(cmdutility.FlagOutput)
			cmdutility.LogError("Failed to encrypt file", err)
			return
		}

		cmdutility.LogColor(cmdutility.Green, "Encryption done. Wrote data to %s", cmdutility.FlagOutput)
	} else {
		profile, profilePassword := cmdutility.GetProfileCommandLine(true)
		if profile == nil {
			return
		}

		rawFilename := path.Base(inputFilename)
		fileAbs, _ := filepath.Abs(inputFilename)
		fileName := filepath.Base(inputFilename)

		var site data.Site
		found := false
		for _, siteresult := range profile.Sites {
			if siteresult.IsFile() {
				if siteresult[data.FileFieldAbsolute] == fileAbs {
					site = siteresult
					found = true
					break
				}
			}
		}
		var fileUuidStr string

		if !found {
			fileUuid, err := uuid.NewRandom()
			if err != nil {
				cmdutility.LogError("Failed to generate file uuid", err)
				return
			}

			fileUuidStr = string(fileUuid.String())

			site = data.NewSiteFile(fileName, fileAbs, fileUuidStr)
		} else {
			fileUuidStr = site[data.FileFieldUUID]
		}

		siteName := rawFilename + "-" + fileUuidStr[len(fileUuidStr)-4:]

		site[data.FileFieldSize] = utils.Kbmbgb(stat.Size())
		site[data.FileFieldUpdate] = time.Now().Local().Format("2006 Jan 2 15:04:05 MST")

		encryptedBytes := bytes.Buffer{}
		encryptStream := bufio.NewWriter(&encryptedBytes)

		cmdutility.LogColor(cmdutility.Green, "Encrypting file %s", inputFilename)
		if err := cipher.StreamEncrypt(inputFile, encryptStream, profilePassword); err != nil {
			os.Remove(cmdutility.FlagOutput)
			cmdutility.LogError("Failed to encrypt file", err)
			return
		}

		encryptStream.Flush()

		site[data.SpecialFieldFileData] = base64.URLEncoding.EncodeToString(encryptedBytes.Bytes())

		profile.AddSite(siteName, site)

		if err := profile.Save(profilePassword); err != nil {
			cmdutility.LogError("Failed saving profile", err)
			return
		}

		cmdutility.LogColor(cmdutility.Green, "Encryption done. Wrote data to site %s", siteName)
	}

	if fileDeleteOriginal {
		cmdutility.LogColor(cmdutility.Green, "Deleting original file %s", inputFilename)
		if err := os.Remove(inputFilename); err != nil {
			cmdutility.LogError("Failed to delete original file", err)
			return
		}
	}
}

func fileDecrypt(cmd *cobra.Command, args []string) {
	inputFilename := args[0]

	if _, err := os.Stat(inputFilename); os.IsNotExist(err) {
		profile, profilePassword := cmdutility.GetProfileCommandLine(true)
		if profile == nil {
			return
		}

		rx, _ := regexp.Compile("[a-fA-F0-9]{4}")
		isUuid := rx.MatchString(inputFilename)

		sites := profile.Sites
		var site data.Site
		siteFound := false
		for siteName, s := range sites {
			if s.IsFile() {
				siteNameParts := strings.Split(siteName, "-")
				if isUuid {
					if siteNameParts[1] == inputFilename {
						site = s
						siteFound = true
						break
					}
				} else {
					if siteNameParts[0] == inputFilename {
						site = s
						siteFound = true
						break
					}
				}
			}
		}

		if siteFound {
			if len(cmdutility.FlagOutput) == 0 {
				cmdutility.FlagOutput = site[data.FileFieldName] + ".decrypted"
			}

			outputFile, err := os.Create(cmdutility.FlagOutput)
			if err != nil {
				cmdutility.LogError("Failed to create output file", err)
				return
			}

			fileBytes := bytes.Buffer{}
			fileBytes.Write(site.GetFileBytes())
			inputReader := bufio.NewReader(&fileBytes)

			cmdutility.LogColor(cmdutility.Green, "Decrypting site file %s", site[data.FileFieldName])
			if err := cipher.StreamDecrypt(inputReader, outputFile, profilePassword); err != nil {
				os.Remove(cmdutility.FlagOutput)
				cmdutility.LogError("Failed to decrypt file", err)
				return
			}

			cmdutility.LogColor(cmdutility.Green, "Decryption done. Wrote data to %s", cmdutility.FlagOutput)
		} else {
			cmdutility.LogError("File doesn't exist", err)
		}

		return
	}

	if len(cmdutility.FlagOutput) == 0 {
		cmdutility.FlagOutput = strings.TrimSuffix(inputFilename, ".enc")
	}

	if cmdutility.FlagOutput == inputFilename {
		cmdutility.LogColor(cmdutility.BoldRed, "Input file name %s is same as output %s", inputFilename, cmdutility.FlagOutput)
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

	outputFile, err := os.Create(cmdutility.FlagOutput)
	if err != nil {
		cmdutility.LogError("Failed to create output file", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Decrypting file %s", inputFilename)
	if err := cipher.StreamDecrypt(inputFile, outputFile, password); err != nil {
		os.Remove(cmdutility.FlagOutput)
		cmdutility.LogError("Failed to decrypt file", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Decryption done. Wrote data to %s", cmdutility.FlagOutput)
}
