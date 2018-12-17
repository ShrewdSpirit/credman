package main

import (
	"os"
	"path"

	"github.com/ShrewdSpirit/credman/utility"

	"github.com/ShrewdSpirit/credman/data"

	"github.com/ShrewdSpirit/credman/cmd"
	homedir "github.com/mitchellh/go-homedir"
)

func checkDataDir() {
	home, _ := homedir.Dir()
	data.DataDir = path.Join(home, ".credman")

	stat, err := os.Stat(data.DataDir)
	if err != nil {
		os.Mkdir(data.DataDir, os.ModePerm)
	} else if !stat.IsDir() {
		os.Remove(data.DataDir)
		os.Mkdir(data.DataDir, os.ModePerm)
	}

	data.ProfilesDir = path.Join(data.DataDir, "profiles")
	os.Mkdir(data.ProfilesDir, os.ModePerm)
}

func main() {
	checkDataDir()

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
