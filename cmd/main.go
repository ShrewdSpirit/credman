package main

import (
	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/data"
	"github.com/spf13/cobra/cobra/cmd"
)

func main() {
	if err := data.LoadConfiguration(); err != nil {
		cmdutility.LogError("Failed to load config", err)
		return
	}

	cmd.Execute()

	if err := data.SaveConfiguration(); err != nil {
		cmdutility.LogError("Failed to save config", err)
		return
	}
}
