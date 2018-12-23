package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

var DataDir string
var Config Configuration

type Configuration struct {
	DefaultProfile string
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
