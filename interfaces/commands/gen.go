package commands

import (
	"fmt"
	"strings"

	"github.com/ShrewdSpirit/credman/interfaces/commands/cmdutility"
	"github.com/ShrewdSpirit/credman/utils"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var genPlen byte
var genPcase string
var genPmix []string
var genCopy bool

var genCmd = &cobra.Command{
	Use:     "gen",
	Aliases: []string{"g"},
	Short:   "Generate password",
	Long:    `Generates a password using credman's password generator`,
	Run: func(cmd *cobra.Command, args []string) {
		var pcase utils.PasswordCase
		pmix := make([]utils.PasswordMix, 1)

		switch strings.ToLower(genPcase) {
		case "both":
			pcase = utils.PasswordCaseBoth
		case "lower":
			pcase = utils.PasswordCaseLower
		case "upper":
			pcase = utils.PasswordCaseUpper
		default:
			fmt.Println("Invalid password generator case")
			return
		}

		for _, mix := range genPmix {
			switch strings.ToLower(mix) {
			case "letter":
				pmix = append(pmix, utils.PasswordMixLetter)
			case "digit":
				pmix = append(pmix, utils.PasswordMixDigit)
			case "punc":
				pmix = append(pmix, utils.PasswordMixPunc)
			case "all":
				pmix = []utils.PasswordMix{utils.PasswordMixAll}
			default:
				fmt.Println("Invalid password generator mix")
				return
			}
		}

		password := utils.GeneratePassword(genPlen, pcase, pmix...)
		if genCopy {
			if err := clipboard.WriteAll(password); err != nil {
				cmdutility.LogError("Failed writing to clipboard", err)
				return
			}
			fmt.Println("Generated password copied to clipboard")
		} else {
			cmdutility.LogColor(cmdutility.HiGreen, "%s", password)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().Uint8VarP(&genPlen, "password-length", "l", 16, "Password generation length")
	genCmd.Flags().StringVarP(&genPcase, "password-case", "c", "both", "Password generation letter case: lower,upper,both")
	genCmd.Flags().StringSliceVarP(&genPmix, "password-mix", "m", []string{"all"}, "Password generation character mix: letter,digit,punc,all")
	genCmd.Flags().BoolVar(&genCopy, "copy", false, "Copy first selected field into clipboard")
}
