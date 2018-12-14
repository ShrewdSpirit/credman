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

func NewSite(siteName string, password string, cmd *cobra.Command) (site *Site) {
	site = &Site{
		Name:     siteName,
		Password: password,
	}

	site.Email, _ = cmd.Flags().GetString("email")
	site.Username, _ = cmd.Flags().GetString("username")
	site.Notes, _ = cmd.Flags().GetString("notes")
	site.SecurityQuestions[0], _ = cmd.Flags().GetString("secq1")
	site.SecurityQuestions[1], _ = cmd.Flags().GetString("secq2")
	site.SecurityQuestions[2], _ = cmd.Flags().GetString("secq3")
	site.SecurityQuestions[3], _ = cmd.Flags().GetString("secq4")
	site.SecurityQuestions[4], _ = cmd.Flags().GetString("secq5")

	return
}
