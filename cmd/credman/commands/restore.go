package commands

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/ShrewdSpirit/credman/cipher"

	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/data"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"rs"},
	Short:   "Profile password restore",
	Run:     restore,
}

var restoreAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add password restore to profile",
	Run:     restoreAdd,
}

var restoreRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r"},
	Short:   "Remove restore from profile",
	Run:     restoreRemove,
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	cmdutility.FlagsAddProfileName(restoreCmd)

	restoreCmd.AddCommand(restoreAddCmd)
	cmdutility.FlagsAddProfileName(restoreAddCmd)

	restoreCmd.AddCommand(restoreRemoveCmd)
	cmdutility.FlagsAddProfileName(restoreRemoveCmd)
}

func restoreReadSecurityQuestions(questionOrders []int) ([]string, []int) {
	answers := make([]string, 0)
	orders := make([]int, 0)

	if questionOrders == nil {
		fmt.Println("You must answer at least 3 questions from the list bellow. Leave any answer empty to ignore the question.")
		for index, q := range data.Config.RestoreQuestions {
			cmdutility.LogColor(cmdutility.Green, " %s - %s", strconv.Itoa(index+1), q)
		}
		fmt.Println("-----------------------")

		reader := bufio.NewReader(os.Stdin)
		for index, question := range data.Config.RestoreQuestions {
			fmt.Printf("[%d/%d] %s ", index+1, len(data.Config.RestoreQuestions), question)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if len(answer) > 0 {
				answers = append(answers, answer)
				orders = append(orders, index)
			}
		}

		if len(answers) < 3 {
			cmdutility.LogColor(cmdutility.BoldRed, "You've answered %s questions which is not enough.", strconv.Itoa(len(answers)))
		}
	} else {
		orders = questionOrders
		fmt.Println("Answer all the questions bellow:")
		reader := bufio.NewReader(os.Stdin)
		for _, order := range questionOrders {
			fmt.Printf("%s ", data.Config.RestoreQuestions[order])
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if len(answer) > 0 {
				answers = append(answers, answer)
			} else {
				fmt.Println("You must answer all questions")
				return nil, nil
			}
		}
	}

	return answers, orders
}

func makeRestoreKey(answers []string, orders []int) []byte {
	sha := sha256.New()
	for _, answer := range answers {
		sha.Write([]byte(answer))
	}
	for _, order := range orders {
		sha.Write([]byte{byte(order % 255)})
	}
	sha.Write([]byte{byte(len(answers) % 255)})

	return sha.Sum(nil)
}

func restore(cmd *cobra.Command, args []string) {
	profile, _ := cmdutility.GetProfileCommandLine(false)
	if profile == nil {
		return
	}

	if profile.Meta.Restore == nil || profile.Meta.RestoreOrder == nil {
		cmdutility.LogColor(cmdutility.Green, "No restore has been set for profile %s", profile.Name)
		return
	}

	answers, orders := restoreReadSecurityQuestions(profile.Meta.RestoreOrder)
	if answers == nil || len(answers) < 3 {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Invalid security questions")
		return
	}

	key := makeRestoreKey(answers, orders)

	profilePassword, err := cipher.BlockDecrypt(profile.Meta.Restore, string(key))
	if err != nil {
		cmdutility.LogError("Failed to decrypt restore", err)
		return
	}

	if err := clipboard.WriteAll(string(profilePassword)); err != nil {
		cmdutility.LogError("Failed to copy to clipboard", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Profile password of %s has been copied to clipboard", profile.Name)
}

func restoreAdd(cmd *cobra.Command, args []string) {
	profile, _ := cmdutility.GetProfileCommandLine(false)
	if profile == nil {
		return
	}

	if profile.Meta.Restore != nil {
		replace, err := cmdutility.YesNoPrompt(fmt.Sprintf("Restore has already setup for profile %s. Are you sure to overwrite restore?", profile.Name))
		if err != nil {
			cmdutility.LogError("Reading input failed", err)
			return
		}

		if !replace {
			return
		}
	}

	answers, orders := restoreReadSecurityQuestions(nil)
	if len(answers) < 3 {
		return
	}

	key := makeRestoreKey(answers, orders)

	profilePassword, err := cmdutility.PasswordPrompt("Profile password")
	if err != nil {
		cmdutility.LogError("Failed reading password", err)
		return
	}

	if profile.Meta.Restore, err = cipher.BlockEncrypt([]byte(profilePassword), string(key)); err != nil {
		cmdutility.LogError("Failed to encrypt restore key", err)
		return
	}

	profile.Meta.RestoreOrder = orders

	if err = profile.SaveRaw(); err != nil {
		cmdutility.LogError("Failed to save profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Restore has been set for profile %s", profile.Name)
}

func restoreRemove(cmd *cobra.Command, args []string) {
	profile, _ := cmdutility.GetProfileCommandLine(false)
	if profile == nil {
		return
	}

	remove, err := cmdutility.YesNoPrompt(fmt.Sprintf("Are you sure to delete restore for profile %s?", profile.Name))
	if err != nil {
		cmdutility.LogError("Reading input failed", err)
		return
	}

	if !remove {
		return
	}

	profile.Meta.Restore = nil
	profile.Meta.RestoreOrder = nil

	if err := profile.SaveRaw(); err != nil {
		cmdutility.LogError("Failed to save profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Restore has been removed from profile %s", profile.Name)
}
