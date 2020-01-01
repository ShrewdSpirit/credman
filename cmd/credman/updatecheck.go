package main

import (
	"fmt"
	"os"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/utils"
)

func checkUpdate() {
	if !utils.DueUpdateCheck() {
		return
	}

	fmt.Println("Checking for updates")
	version, err := utils.CheckNewVersion()
	if err != nil {
		cmdutility.LogError("Failed to check new version", err)
		return
	}

	if version == data.Version {
		fmt.Println("No updates available")
		return
	}

	doUpdate, err := cmdutility.YesNoPrompt(fmt.Sprintf("New version %s available. Get update?", version))
	if err != nil {
		cmdutility.LogError("Reading input failed", err)
	}

	if !doUpdate {
		return
	}

	fmt.Println("Downloading update ...")
	if err := utils.GetUpdate(); err != nil {
		cmdutility.LogError("Failed getting new update", err)
		return
	}

	fmt.Println("Download finished. Installing")
	if err := utils.InstallUpdate(); err != nil {
		cmdutility.LogError("Failed installing update", err)
		return
	}

	os.Exit(0)
}
