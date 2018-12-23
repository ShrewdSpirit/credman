package cmdutility

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/utility"
	"github.com/spf13/cobra"
)

var FlagProfileName string

var FlagPgen bool
var FlagPlen byte
var FlagPcase string
var FlagPmix []string

func FlagsAddProfileName(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&FlagProfileName, "profile", "p", "", "Profile to use")
}

func FlagsAddPasswordOptions(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&FlagPgen, "password-generate", "g", false, "Uses password generation")
	cmd.Flags().Uint8VarP(&FlagPlen, "password-length", "l", 16, "Password generation length")
	cmd.Flags().StringVarP(&FlagPcase, "password-case", "c", "both", "Password generation letter case: lower,upper,both")
	cmd.Flags().StringSliceVarP(&FlagPmix, "password-mix", "m", []string{"all"}, "Password generation character mix: letter,digit,punc,all")
}

func GetProfileCommandLine() (*data.Profile, string) {
	profileName := data.Config.DefaultProfile
	if len(FlagProfileName) != 0 {
		profileName = FlagProfileName
	}

	if len(profileName) == 0 {
		fmt.Println("Create a profile first!")
		return nil, ""
	}

	if !data.ProfileExists(profileName) {
		utility.LogColor(utility.BoldHiYellow, "Profile %s doesnt exist.", profileName)
		return nil, ""
	}

	utility.LogColor(utility.Green, "Using profile %s", profileName)

	password, err := utility.PasswordPrompt("Profile password")
	if err != nil {
		utility.LogError("Failed reading password", err)
		return nil, ""
	}

	profile, err := data.LoadProfile(profileName, password)
	if err != nil {
		utility.LogError("Failed loading profile", err)
		return nil, ""
	}

	return profile, password
}

func ParsePasswordGenerationFlags(prompt string) (string, error) {
	if FlagPgen {
		var pcase utility.PasswordCase
		pmix := make([]utility.PasswordMix, 1)

		switch strings.ToLower(FlagPcase) {
		case "both":
			pcase = utility.PasswordCaseBoth
		case "lower":
			pcase = utility.PasswordCaseLower
		case "upper":
			pcase = utility.PasswordCaseUpper
		default:
			return "", errors.New("Invalid password generator case")
		}

		for _, mix := range FlagPmix {
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
				return "", errors.New("Invalid password generator mix")
			}
		}

		password := utility.GeneratePassword(FlagPlen, pcase, pmix...)
		return password, nil
	}

	password, err := utility.NewPasswordPrompt(prompt)
	if err != nil {
		return "", err
	}

	return password, nil
}
