package framework

import (
	"io/ioutil"

	config "github.com/multiverse-os/starshipyard/framework/config"
	yaml "gopkg.in/yaml.v2"
)

func DefaultConfig() *config.Config {
	return &config.Config{
		AppName:     "starship",
		Description: "a web application framework with a focus on security and heavily inspired by rails",
		Keywords:    []string{"web", "framework", "example", "golang"},
		// may not be necessary since its false by default...
		MaintainanceMode: false,
		Announcement:     "We are down for maintainance, and we will be back shortly.",
		SessionsDisabled: false,
		Environment:      "development",
		Pid:              "tmp/pids/starship.pid",
		Address:          "localhost",
		Port:             3000,
		Debug:            true,
	}
}

func LoadConfig(path string) (config *config.Config, err error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
