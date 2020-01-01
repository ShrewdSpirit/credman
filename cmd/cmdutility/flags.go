package cmdutility

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/utils"
	"github.com/spf13/cobra"
)

var FlagProfileName string

var FlagPgen bool
var FlagPlen byte
var FlagPcase string
var FlagPmix []string
var FlagOutput string

func FlagsAddProfileName(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&FlagProfileName, "profile", "p", "", "Profile to use")
}

func FlagsAddPasswordOptions(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&FlagPgen, "password-generate", "g", false, "Uses password generation")
	cmd.Flags().Uint8VarP(&FlagPlen, "password-length", "l", 16, "Password generation length")
	cmd.Flags().StringVarP(&FlagPcase, "password-case", "c", "both", "Password generation letter case: lower,upper,both")
	cmd.Flags().StringSliceVarP(&FlagPmix, "password-mix", "m", []string{"all"}, "Password generation character mix: letter,digit,punc,all")
}

func FlagsAddOutput(cmd *cobra.Command, desc string) {
	cmd.Flags().StringVarP(&FlagOutput, "output", "o", "", desc)
}

func GetProfileCommandLine(readPassword bool) (*data.Profile, string) {
	profileName := data.Config.DefaultProfile
	if len(FlagProfileName) != 0 {
		profileName = FlagProfileName
	}

	if len(profileName) == 0 {
		fmt.Println("Create a profile first!")
		return nil, ""
	}

	if !data.ProfileExists(profileName) {
		LogColor(BoldHiYellow, "Profile %s doesnt exist.", profileName)
		return nil, ""
	}

	LogColor(Green, "Using profile %s", profileName)

	if !readPassword {
		profile, err := data.LoadProfileRaw(profileName)
		if err != nil {
			LogError("Failed loading profile", err)
			return nil, ""
		}

		return profile, ""
	}

	password, err := PasswordPrompt("Profile password")
	if err != nil {
		LogError("Failed reading password", err)
		return nil, ""
	}

	profile, err := data.LoadProfile(profileName, password)
	if err != nil {
		LogError("Failed loading profile", err)
		return nil, ""
	}

	return profile, password
}

func ParsePasswordGenerationFlags(prompt string) (string, error) {
	if FlagPgen {
		var pcase utils.PasswordCase
		pmix := make([]utils.PasswordMix, 1)

		switch strings.ToLower(FlagPcase) {
		case "both", "b":
			pcase = utils.PasswordCaseBoth
		case "lower", "l":
			pcase = utils.PasswordCaseLower
		case "upper", "u":
			pcase = utils.PasswordCaseUpper
		default:
			return "", errors.New("Invalid password generator case")
		}

		for _, mix := range FlagPmix {
			switch strings.ToLower(mix) {
			case "letter", "l":
				pmix = append(pmix, utils.PasswordMixLetter)
			case "digit", "d":
				pmix = append(pmix, utils.PasswordMixDigit)
			case "punc", "p":
				pmix = append(pmix, utils.PasswordMixPunc)
			case "all", "a":
				pmix = []utils.PasswordMix{utils.PasswordMixAll}
			default:
				return "", errors.New("Invalid password generator mix")
			}
		}

		password := utils.GeneratePassword(FlagPlen, pcase, pmix...)
		return password, nil
	}

	password, err := NewPasswordPrompt(prompt)
	if err != nil {
		return "", err
	}

	return password, nil
}
