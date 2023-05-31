package helpers

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

func GetAppRootURL() string {
	//return "http://localhost:8080"
	return "https://api.getvigil.io"
}

func GetConfigPath() (configPath string, err error) {
	if hp, e := homedir.Dir(); e != nil {
		err = e
		return
	} else {
		configPath = path.Join(hp, ".vigil", "config.yaml")
		return
	}
}

func LoadConfig() (config Config, err error) {
	if os.Getenv("VIGIL_USERNAME") != "" && os.Getenv("VIGIL_TOKEN") != "" {
		config = Config{
			Username: os.Getenv("VIGIL_USERNAME"),
			Token:    os.Getenv("VIGIL_TOKEN"),
		}
		return
	}
	if configPath, e := GetConfigPath(); e != nil {
		err = e
		return
	} else if yamlFile, e := ioutil.ReadFile(configPath); e != nil {
		err = e
	} else if e = yaml.Unmarshal(yamlFile, &config); e != nil {
		err = e
	}
	if os.Getenv("VIGIL_USERNAME") != "" {
		config.Username = os.Getenv("VIGIL_USERNAME")
	}
	if os.Getenv("VIGIL_TOKEN") != "" {
		config.Token = os.Getenv("VIGIL_TOKEN")
	}
	return
}
