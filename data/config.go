package data

import (
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	Version   = "N/A"
	GitCommit = "N/A"
)

var DataDir string
var ProfilesDir string
var Config Configuration

type Configuration struct {
	DefaultProfile     string
	AutoUpdateInterval int `json:"AutoUpdateDays"` // 0 means no autoupdate
	RestoreQuestions   []string
	WebInterfacePort   int   `json:"webinterface_port"`
	LastUpdateCheck    int64 `json:"p_uchk"`
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
	Config.AutoUpdateInterval = 1
	Config.WebInterfacePort = 14201
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
	setDefaultConfig()

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
	data, err := json.MarshalIndent(&Config, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(DataDir, "config.json"), data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
