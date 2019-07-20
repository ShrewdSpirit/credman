package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/ShrewdSpirit/credman/interfaces/browser"
	"github.com/ShrewdSpirit/credman/utils"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "credman",
	Short: "Simple yet powerful credential/password manager",
	Long: `Simple yet powerful credential/password manager with remote sync.

Use 'credman help password' to see how to use password options.
Use 'credman help fields to see what fields you can set for a site.`,
	Run: func(cmd *cobra.Command, args []string) {
		browser.Run()
	},
}

var rootClsClip = &cobra.Command{
	Use:    "clsclip",
	Hidden: true,
	Run:    clsClip,
}

func init() {
	rootCmd.AddCommand(rootClsClip)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func clsClip(cmd *cobra.Command, args []string) {
	time.Sleep(10 * time.Second)
	clipboard.WriteAll("")
	utils.RemovePidFile()
}
