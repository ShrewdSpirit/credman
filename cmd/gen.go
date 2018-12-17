package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/ShrewdSpirit/credman/utility"
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
		var pcase utility.PasswordCase
		pmix := make([]utility.PasswordMix, 1)

		switch strings.ToLower(genPcase) {
		case "both":
			pcase = utility.PasswordCaseBoth
		case "lower":
			pcase = utility.PasswordCaseLower
		case "upper":
			pcase = utility.PasswordCaseUpper
		default:
			fmt.Println("Invalid password generator case")
			return
		}

		for _, mix := range genPmix {
			switch strings.ToLower(mix) {
			case "letter":
				pmix = append(pmix, utility.PasswordMixLetter)
			case "digit":
				pmix = append(pmix, utility.PasswordMixDigit)
			case "punc":
				pmix = append(pmix, utility.PasswordMixPunc)
			case "all":
				pmix = []utility.PasswordMix{utility.PasswordMixAll}
			default:
				fmt.Println("Invalid password generator mix")
				return
			}
		}

		password := utility.GeneratePassword(genPlen, pcase, pmix...)
		if genCopy {
			if err := clipboard.WriteAll(password); err != nil {
				utility.LogError("Failed writing to clipboard", err)
				return
			}
			fmt.Println("Generated password copied to clipboard")
		} else {
			utility.LogColor(utility.HiGreen, "%s", password)
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
