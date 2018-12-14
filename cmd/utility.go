package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cobra"
)

var ProfilesDir string

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
	return f == "letter" || f == "digit" || f == "both"
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

func UsesPasswordGenerator(cmd *cobra.Command) bool {
	f, err := cmd.Flags().GetBool("pgen")
	if err != nil {
		return false
	}
	return f
}

func PromptPassword(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	password := strings.TrimSpace(string(passwordBytes))

	if len(password) >= 32 {
		return "", errors.New("Password can't be longer than 32 characters")
	}

	return password, nil
}

func GetPassword(cmd *cobra.Command, prompt string) (string, error) {
	if cmd != nil && UsesPasswordGenerator(cmd) {
		// generate
		if !PasswordValidatePcase(cmd) {
			pcase, _ := cmd.Flags().GetString("pcase")
			return "", errors.New(fmt.Sprintf("Invalid password case: %s", pcase))
		}
		if !PasswordValidatePmix(cmd) {
			pmix, _ := cmd.Flags().GetString("pmix")
			return "", errors.New(fmt.Sprintf("Invalid password mix: %s", pmix))
		}
		plen, err := cmd.Flags().GetUint8("plen")
		if err != nil {
			return "", err
		}
		if plen >= 32 {
			return "", errors.New("Password can't be longer than 32 characters")
		}

		pcase := PasswordGetPcase(cmd)
		pmix := PasswordGetPmix(cmd)
		pw := GeneratePassword(plen, pcase, pmix)
		return pw, nil
	}

	if prompt == "" {
		prompt = "Password"
	}
	password, err := PromptPassword(prompt)
	if err != nil {
		return "", err
	}
	repeatPassword, err := PromptPassword("Repeat password")
	if err != nil {
		return "", err
	}

	if password != repeatPassword {
		return "", errors.New("Passwords doesn't match")
	}

	return password, nil
}

var lowerLetters = "abcdefghijklmnopqrstuv"
var upperLetters = "ABCDEFGHIJKLMNOPQRSTUV"
var digitLetters = "0123456789"

func GeneratePassword(plen byte, pcase PasswordCase, pmix PasswordMix) string {
	gen := func(dict string) string {
		src := rand.NewSource(time.Now().UnixNano())
		b := []byte(dict)
		l := int64(len(b))
		buf := bytes.Buffer{}
		i := byte(0)
		for ; i < plen; i++ {
			buf.WriteByte(b[src.Int63()%l])
		}
		return buf.String()
	}

	switch pmix {
	case PasswordMixLetter:
		switch pcase {
		case PasswordCaseLower:
			return gen(lowerLetters)
		case PasswordCaseUpper:
			return gen(upperLetters)
		case PasswordCaseBoth:
			return gen(lowerLetters + upperLetters)
		}
	case PasswordMixBoth:
		switch pcase {
		case PasswordCaseLower:
			return gen(digitLetters + lowerLetters)
		case PasswordCaseUpper:
			return gen(digitLetters + upperLetters)
		case PasswordCaseBoth:
			return gen(digitLetters + lowerLetters + upperLetters)
		}
	case PasswordMixDigit:
		return gen(digitLetters)
	}

	return ""
}

func ProfileAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("profile", "p", "", "Specifies which profile to use")
}

func SiteAddFieldFlags(cmd *cobra.Command, isBool bool) {
	if isBool {
		cmd.Flags().Bool("password", false, "Gets password field")
		cmd.Flags().Bool("email", false, "Gets email field")
		cmd.Flags().Bool("username", false, "Gets username field")
		cmd.Flags().Bool("notes", false, "Gets notes field")
		cmd.Flags().Bool("secq1", false, "Gets first security question")
		cmd.Flags().Bool("secq2", false, "Gets second security question")
		cmd.Flags().Bool("secq3", false, "Gets third security question")
		cmd.Flags().Bool("secq4", false, "Gets forth security question")
		cmd.Flags().Bool("secq5", false, "Gets fifth security question")
	} else {
		cmd.Flags().String("email", "", "Sets email field")
		cmd.Flags().String("username", "", "Sets username field")
		cmd.Flags().String("notes", "", "Sets notes field")
		cmd.Flags().String("secq1", "", "Sets first security question")
		cmd.Flags().String("secq2", "", "Sets second security question")
		cmd.Flags().String("secq3", "", "Sets third security question")
		cmd.Flags().String("secq4", "", "Sets forth security question")
		cmd.Flags().String("secq5", "", "Sets fifth security question")
	}
}
