package cmd

import (
	"github.com/spf13/cobra"
)

var fieldsCmd = &cobra.Command{
	Use:   "fields",
	Short: "Fields help",
	Long: `Fields are the details of each site. A site always have a password field that stores the generated or user input password for that site.
You can set any field for a site. For setting custom fields you must use -f/--field option with this format:
-f="field key"="field value" or -f="field key"="field value" -f=another=field
The quotes are optional if there's no space character in the name.

Also site get's field option takes a list of fields instead of key value pairs: -f=email,password,etc`,
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(fieldsCmd)
}
