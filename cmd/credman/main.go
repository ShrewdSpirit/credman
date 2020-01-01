package main

import (
	"os"
	"path"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/cmd"
	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/utils"
)

func main() {
	if err := data.LoadConfiguration(); err != nil {
		cmdutility.LogError("Failed to load config", err)
		return
	}

	if result, _ := utils.IsSlaveProcess(); !result {
		if _, err := os.Stat(path.Join(data.DataDir, "update")); err == nil {
			os.Remove(path.Join(data.DataDir, "update"))
		}

		if len(os.Args) > 1 {
			checkUpdate()
		}
	}

	commands.Execute()

	if err := data.SaveConfiguration(); err != nil {
		cmdutility.LogError("Failed to save config", err)
		return
	}
}
