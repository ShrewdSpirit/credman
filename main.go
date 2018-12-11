package main

import (
	"os"
	"path"

	"github.com/ShrewdSpirit/credman/cmd"
	"github.com/ShrewdSpirit/credman/config"
	homedir "github.com/mitchellh/go-homedir"
)

func checkDataDir() {
	home, _ := homedir.Dir()
	dataDir := path.Join(home, ".credman")
	stat, err := os.Stat(dataDir)
	if err != nil {
		os.Mkdir(dataDir, os.ModePerm)
	} else if !stat.IsDir() {
		os.Remove(dataDir)
		os.Mkdir(dataDir, os.ModePerm)
	}
	cmd.ProfilesDir = path.Join(dataDir, "profiles")
	os.Mkdir(cmd.ProfilesDir, os.ModePerm)
}

func main() {
	checkDataDir()
	config.LoadConfig()
	cmd.Execute()
	config.SaveConfig()
}
