package framework

import (
	"fmt"
	"io/ioutil"

	config "github.com/multiverse-os/starshipyard/framework/config"
	yaml "gopkg.in/yaml.v2"
)

func DefaultConfig() config.Settings {
	return config.Settings{
		Address: "localhost",
		Port:    3000,
		Pid:     "tmp/pids/starship.pid",
		MaintainanceMode: config.Maintainance{
			Enabled:      false,
			Announcement: "We are down for maintainance, and we will be back shortly.",
			UserSessions: false,
		},
	}
}

func LoadConfig(path string) (config config.Settings, err error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func ValidateConfig(config config.Settings) config.Settings {
	// TODO: Need more validations for all the individual fields
	// TODO: Port needs to only support actual ports 1 - ~65000
	if len(config.Pid) == 0 {
		config.Pid = "tmp/pids/starship.pid"
		fmt.Println("config.Pid set to because it was blank:", config.Pid)
	}
	if len(config.DataDirectory) == 0 {
		config.DataDirectory = "data"
	}
	return config
}
