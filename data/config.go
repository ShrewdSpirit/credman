package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

var DataDir string
var Version string = "0.8.0"
var Config Configuration

type Configuration struct {
	DefaultProfile     string
	AutoUpdateInterval int // 0 means no autoupdate
	RestoreQuestions   []string
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

func setDefaultConfig() {
	Config.AutoUpdateInterval = 0

	Config.RestoreQuestions = []string{
		"What's the name of your first high school?",
		"What are the last 5 digits of your ID?",
		"What's your favorite animal?",
		"What's your favorite movie?",
		"How many siblings do you have?",
		"In what city or town did your mother and father meet?",
		"What's your favorite hobby?",
		"What are the last 4 digits of your credit card?",
		"In what city was your first job?",
		"What's your first nickname?",
		"What's your favorite restaurant?",
		"What is the first name of the boy or girl that you first kissed?",
		"What is your preferred musical genre?",
	}
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

	if Config.RestoreQuestions == nil || len(Config.RestoreQuestions) < 3 {
		setDefaultConfig()
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
