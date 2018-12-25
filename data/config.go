package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

var DataDir string
var Version string = "0.3.2"
var Config Configuration

type Configuration struct {
	DefaultProfile string
}

func init() {
	home, _ := homedir.Dir()
	DataDir = path.Join(home, ".credman")

	stat, err := os.Stat(DataDir)
	if err != nil {
		os.Mkdir(DataDir, os.ModePerm)
	} else if !stat.IsDir() {
		os.Remove(DataDir)
		os.Mkdir(DataDir, os.ModePerm)
	}

	ProfilesDir = path.Join(DataDir, "profiles")
	os.Mkdir(ProfilesDir, os.ModePerm)
}

func LoadConfiguration() error {
	configFilePath := path.Join(DataDir, "config.json")
	_, err := os.Stat(configFilePath)
	if err != nil {
		return nil
	}

	configData, err := ioutil.ReadFile(path.Join(DataDir, "config.json"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(configData, &Config)
	if err != nil {
		return err
	}

	return nil
}

func SaveConfiguration() error {
	data, err := json.MarshalIndent(&Config, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(DataDir, "config.json"), data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
