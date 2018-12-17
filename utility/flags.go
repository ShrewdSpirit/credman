package utility

import (
	"github.com/spf13/cobra"
)

var FlagProfileName string

func FlagsAddProfileName(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&FlagProfileName, "profile", "p", "", "Profile to use")
}
