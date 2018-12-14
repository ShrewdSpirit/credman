package config

import (
	"github.com/spf13/cobra"
)

type Site struct {
	Name   string
	Fields map[string]string
}

func NewSite(siteName string, password string, cmd *cobra.Command) (site *Site) {
	site = &Site{
		Name:   siteName,
		Fields: make(map[string]string),
	}

	site.Fields["email"], _ = cmd.Flags().GetString("email")
	site.Fields["username"], _ = cmd.Flags().GetString("username")
	site.Fields["password"] = password
	fields, _ := cmd.Flags().GetStringToString("fields")
	for fkey, fvalue := range fields {
		site.Fields[fkey] = fvalue
	}

	return
}
