package main

import (
	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/cmd/commands"
	"github.com/ShrewdSpirit/credman/data"
)

func main() {
	if err := data.LoadConfiguration(); err != nil {
		cmdutility.LogError("Failed to load config", err)
		return
	}

	commands.Execute()

	if err := data.SaveConfiguration(); err != nil {
		cmdutility.LogError("Failed to save config", err)
		return
	}
}
