package cmd

import "github.com/spf13/cobra"

type PasswordCase int
type PasswordMix int

const (
	PasswordCaseBoth  PasswordCase = 0
	PasswordCaseLower PasswordCase = 1
	PasswordCaseUpper PasswordCase = 2
)

const (
	PasswordMixBoth   PasswordMix = 0
	PasswordMixLetter PasswordMix = 1
	PasswordMixDigit  PasswordMix = 2
)

func PasswordAddFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("pgen", false, "Use password generator")
	cmd.Flags().Uint8("plen", 12, "Password generator's password length")
	cmd.Flags().String("pcase", "both", "Password generator's character case: lower/upper/both")
	cmd.Flags().String("pmix", "both", "Password generator's characters type: letter/digit/both")
}

func PasswordValidatePcase(cmd *cobra.Command) bool {
	f, err := cmd.Flags().GetString("pcase")
	if err != nil {
		return false
	}
	return f == "lower" || f == "upper" || f == "both"
}

func PasswordGetPcase(cmd *cobra.Command) PasswordCase {
	f, _ := cmd.Flags().GetString("pcase")
	switch f {
	case "both":
		return PasswordCaseBoth
	case "lower":
		return PasswordCaseLower
	case "upper":
		return PasswordCaseUpper
	default:
		return PasswordCaseBoth
	}
}

func PasswordValidatePmix(cmd *cobra.Command) bool {
	f, err := cmd.Flags().GetString("pmix")
	if err != nil {
		return false
	}
	return f == "letter" || f == "digit" || f == "mix"
}

func PasswordGetPmix(cmd *cobra.Command) PasswordMix {
	f, _ := cmd.Flags().GetString("pmix")
	switch f {
	case "both":
		return PasswordMixBoth
	case "letter":
		return PasswordMixLetter
	case "digit":
		return PasswordMixDigit
	default:
		return PasswordMixBoth
	}
}

func UsePasswordGenerator(cmd *cobra.Command) bool {
	f, err := cmd.Flags().GetBool("pgen")
	if err != nil {
		return false
	}
	return f
}

func ProfileAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("profile", "p", "default", "Specifies which profile to use")
}
