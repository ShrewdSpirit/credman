package main

import (
	"github.com/ShrewdSpirit/credman/utility"

	"github.com/ShrewdSpirit/credman/data"

	"github.com/ShrewdSpirit/credman/cmd"
)

func main() {
	if err := data.LoadConfiguration(); err != nil {
		utility.LogError("Failed to load config", err)
		return
	}

	cmd.Execute()

	if err := data.SaveConfiguration(); err != nil {
		utility.LogError("Failed to save config", err)
		return
	}
}
