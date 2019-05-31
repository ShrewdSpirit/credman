package main

import (
	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/cmd/commands"
	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/utility"
)

func main() {
	if err := data.LoadConfiguration(); err != nil {
		cmdutility.LogError("Failed to load config", err)
		return
	}

	if result, _ := utility.IsSlaveProcess(); !result {
		checkUpdate()
	}

	commands.Execute()

	if err := data.SaveConfiguration(); err != nil {
		cmdutility.LogError("Failed to save config", err)
		return
	}
}
