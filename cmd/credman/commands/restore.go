package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ShrewdSpirit/credman/cmd/credman/cmdutility"
	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/management"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"rs"},
	Short:   "Profile password restore",
	Run: func(cmd *cobra.Command, args []string) {
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
			return
		}

		management.RestoreData{
			Answers: answers,
			Orders:  orders,
			Profile: profile,
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.RestoreStepDecrypting:
						fmt.Println("Restoring")
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Profile password of %s has been copied to clipboard", profile.Name)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.RestoreStepDecrypting:
						cmdutility.LogError("Failed to decrypt restore", err)
					case management.RestoreStepClipboardPassword:
						cmdutility.LogError("Failed to copy to clipboard", err)
					}
				},
			},
		}.Restore()
	},
}

var restoreAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add password restore to profile",
	Run: func(cmd *cobra.Command, args []string) {
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

		management.RestoreData{
			Answers: answers,
			Orders:  orders,
			Profile: profile,
			PasswordReader: func(step management.ManagementStep) string {
				password, err := cmdutility.PasswordPrompt("Profile password")
				if err != nil {
					cmdutility.LogError("Failed reading password", err)
					return ""
				}
				return password
			},
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.RestoreStepEncrypting:
						fmt.Println("Adding restore")
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Restore has been set for profile %s", profile.Name)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.RestoreStepEncrypting:
						cmdutility.LogError("Failed to encrypt restore", err)
					case management.ProfileStepSaving:
						cmdutility.LogError("Failed to save profile", err)
					}
				},
			},
		}.Add()
	},
}

var restoreRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r"},
	Short:   "Remove restore from profile",
	Run: func(cmd *cobra.Command, args []string) {
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

		management.RestoreData{
			Profile: profile,
			ManagementData: management.ManagementData{
				OnStep: func(step management.ManagementStep) {
					switch step {
					case management.StepDone:
						cmdutility.LogColor(cmdutility.Green, "Restore has been removed from profile %s", profile.Name)
					}
				},
				OnError: func(step management.ManagementStep, err error) {
					switch step {
					case management.ProfileStepSaving:
						cmdutility.LogError("failed to save profile", err)
					}
				},
			},
		}.Remove()
	},
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
