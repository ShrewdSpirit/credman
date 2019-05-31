package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ShrewdSpirit/credman/utility"

	"github.com/ShrewdSpirit/credman/cmd/cmdutility"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/spf13/cobra"
)

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update setting and check",
}

var cmdUpdateCheck = &cobra.Command{
	Use:   "check",
	Short: "Checks for newer version",
	Run:   updateCheck,
}

var cmdUpdateSet = &cobra.Command{
	Use:   "set",
	Short: "Sets auto update interval in days",
	Args:  cobra.ExactArgs(1),
	Run:   updateSet,
}

var cmdUpdateInstall = &cobra.Command{
	Use:    "install",
	Hidden: true,
	Run:    updateInstall,
}

func init() {
	rootCmd.AddCommand(cmdUpdate)

	cmdUpdate.AddCommand(cmdUpdateCheck)

	cmdUpdate.AddCommand(cmdUpdateSet)

	cmdUpdate.AddCommand(cmdUpdateInstall)
}

func updateCheck(cmd *cobra.Command, args []string) {
	fmt.Println("Checking for updates")
	version, err := utility.CheckNewVersion()
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
		fmt.Println("credman is up to date!")
		return
	}

	fmt.Println("Downloading update ...")
	if err := utility.GetUpdate(); err != nil {
		cmdutility.LogError("Failed getting new update", err)
		return
	}

	fmt.Println("Download finished. Installing")
	if err := utility.InstallUpdate(); err != nil {
		cmdutility.LogError("Failed installing update", err)
		return
	}

	os.Exit(0)
}

func updateSet(cmd *cobra.Command, args []string) {
	interval, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		cmdutility.LogError("Failed parsing interval", err)
		return
	}

	if interval == 0 {
		fmt.Println("Auto update disabled")
	} else {
		data.Config.AutoUpdateInterval = int(interval)
		cmdutility.LogColor(cmdutility.HiGreen, "Auto update check interval has been set to %s days", args[0])
	}
}

func updateInstall(cmd *cobra.Command, args []string) {
	utility.RemovePidFile()
	exepathFilename := path.Join(data.DataDir, "origpath")

	time.Sleep(1 * time.Second)

	// TODO: log the error to file

	if _, err := os.Stat(exepathFilename); os.IsNotExist(err) {
		cmdutility.LogError("[credman update] Invalid file", err)
		return
	}

	origExePathBytes, err := ioutil.ReadFile(exepathFilename)
	if err != nil {
		cmdutility.LogError("[credman update] Failed reading original path", err)
		return
	}

	currentExe, err := os.Executable()
	if err != nil {
		cmdutility.LogError("[credman update] Failed getting own executable path", err)
		return
	}

	updateFile, err := os.Open(currentExe)
	if err != nil {
		cmdutility.LogError("[credman update] Failed opening current executable", err)
		return
	}
	defer updateFile.Close()

	originalFile, err := os.Create(string(origExePathBytes))
	if err != nil {
		cmdutility.LogError("[credman update] Failed recreating old executable", err)
		return
	}
	defer originalFile.Close()

	if _, err = io.Copy(originalFile, updateFile); err != nil {
		cmdutility.LogError("[credman update] Failed copying update", err)
		return
	}

	os.Remove(exepathFilename)
	os.Remove(path.Join(data.DataDir, "update.tar.gz"))
}
