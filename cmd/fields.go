package cmd

import (
	"github.com/spf13/cobra"
)

var fieldsCmd = &cobra.Command{
	Use:   "fields",
	Short: "Fields help",
	Long: `Fields are the details of each site. A site always have a password field that stores the generated or user input password for that site.
You can set any field for a site. There are 3 default fields that you can use the helper options, email, username and password.
For using any other custom field you must use -f/--fields option with this format:
-f="field key"="field value" or -f="field key"="field value" -f=another=field
The quotes are optional if there's no space character in the name.`,
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(fieldsCmd)
}
