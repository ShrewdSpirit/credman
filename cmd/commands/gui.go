package commands

import (
	"github.com/ShrewdSpirit/credman/gui"
	"github.com/spf13/cobra"
)

var cmdGui = &cobra.Command{
	Short: "Opens GUI",
	Use:   "gui",
	Run: func(cmd *cobra.Command, args []string) {
		gui.Open()
	},
}

func init() {
	rootCmd.AddCommand(cmdGui)
}
