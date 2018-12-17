package utility

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

func PasswordPrompt(prompt string) (string, error) {
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

func NewPasswordPrompt(prompt string) (string, error) {
	if prompt == "" {
		prompt = "Password"
	}
	password, err := PasswordPrompt(prompt)
	if err != nil {
		return "", err
	}
	repeatPassword, err := PasswordPrompt("Repeat password")
	if err != nil {
		return "", err
	}

	if password != repeatPassword {
		return "", errors.New("Passwords doesn't match")
	}

	return password, nil
}

func YesNoPrompt(message string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(color.Output, message+" %s ", BoldHiYellow.SprintFunc()("[y/n]"))
	c, err := reader.ReadByte()
	if err != nil {
		return false, err
	}
	val := strings.ToLower(string(c))
	return val == "y", nil
}
