package config

import (
	"github.com/spf13/cobra"
)

type Site struct {
	Name              string
	Password          string
	Email             string
	Username          string
	Notes             string
	SecurityQuestions [5]string
}

func NewSite(siteName string, password string, cmd *cobra.Command) (site *Site, err error) {
	site = &Site{
		Name:     siteName,
		Password: password,
	}

	site.Email, err = cmd.Flags().GetString("email")
	if err != nil {
		return
	}

	site.Username, err = cmd.Flags().GetString("username")
	if err != nil {
		return
	}

	site.Notes, err = cmd.Flags().GetString("notes")
	if err != nil {
		return
	}

	site.SecurityQuestions[0], err = cmd.Flags().GetString("secq1")
	if err != nil {
		return
	}

	site.SecurityQuestions[1], err = cmd.Flags().GetString("secq2")
	if err != nil {
		return
	}

	site.SecurityQuestions[2], err = cmd.Flags().GetString("secq3")
	if err != nil {
		return
	}

	site.SecurityQuestions[3], err = cmd.Flags().GetString("secq4")
	if err != nil {
		return
	}

	site.SecurityQuestions[4], err = cmd.Flags().GetString("secq5")
	if err != nil {
		return
	}

	return
}
