package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Format uint8

const (
	YAML Format = iota
	JSON
	TOML
	INI
)

type File struct {
	Path   string
	Format string
}

type Environment int

const (
	Development Environment = iota
	Testing
	Production
)

func (self Environment) String() string {
	switch self {
	case Testing:
		return "testing"
	case Production:
		return "production"
	default:
		return "development"
	}
}

type Maintainance struct {
	Enabled      bool   `yaml:"enabled"`
	UserSessions bool   `yaml:"user_sessions"`
	Announcement string `yaml:"announcement"`
}

// TODO: Would prefer it to all be under app: like seen in rails,  and this can
// be done using our own marshal and unmarshal funcitons would should be done
// regardless for maximum control
// TODO: Address/Port should be handled in an nginx like configuration since
// this application framework is meant to be able to handle reverse proxy,
// multiple hosts/domains
type Settings struct {
	Environment      string       `yaml:"environment"`
	Address          string       `yaml:"address"`
	Port             int          `yaml:"port"`
	Pid              string       `yaml:"pid"`
	DataDirectory    string       `yaml:"data"`
	MaintainanceMode Maintainance `yaml:"maintainance"`
}

// TODO: Add ability to update a setting. Write a default settings file.
// TODO: Separate out application specific settings from config library logic so
// we can re-use this code.

func LoadConfig(path string) (config *Settings, err error) {
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

func (self *Settings) Save(path string) error {
	configPath, _ := filepath.Split(path)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	} else {
		yamlData, err := yaml.Marshal(&self)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(path, yamlData, 0600)
	}
}

// Initialize - First run config folder structure and file
// initialization using default config, unless otherwise
// specified using flags.
func (self *Settings) InitializeDefault(path string) error {
	configPath, _ := filepath.Split(path)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, 0700)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := self.Save(path)
		if err != nil {
			return nil
		}
	}
	return nil
}
