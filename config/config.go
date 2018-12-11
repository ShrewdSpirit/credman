package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

var AppConfig *Config = &Config{}

type Config struct {
	DefaultProfile string
}

func LoadConfig() {
	home, _ := homedir.Dir()
	configDir := path.Join(home, ".credman")
	configData, err := ioutil.ReadFile(path.Join(configDir, "config.json"))
	if err != nil {
		return
	}

	err = json.Unmarshal(configData, AppConfig)
	if err != nil {
		panic(err)
	}
}

func SaveConfig() {
	home, _ := homedir.Dir()
	configDir := path.Join(home, ".credman")
	data, err := json.MarshalIndent(AppConfig, "", "\t")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(configDir, "config.json"), data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
