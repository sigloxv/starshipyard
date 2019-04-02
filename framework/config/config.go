package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Environment int

const (
	Development Environment = iota
	Testing
	Staging
	Production
	Maintainance
)

func (self Environment) String() string {
	switch self {
	case Testing:
		return "testing"
	case Staging:
		return "staging"
	case Production:
		return "production"
	case Maintainance:
		return "maintainance"
	default:
		return "development"
	}
}

// TODO: Not a fan of local de
type Config struct {
	AppName            string   `yaml:"name"`
	Description        string   `yaml:"description"`
	Keywords           []string `yaml:"keywords"`
	MaintainanceMode   bool     `yaml:"maintainance"`
	Announcement       string   `yaml:"announcement"`
	SessionsDisabled   bool     `yaml:"sessions"`
	Environment        string   `yaml:"environment"`
	ConfigDirectory    string   `yaml:"config"`
	TemporaryDirectory string   `yaml:"temp"`
	DataDirectory      string   `yaml:"data"`
	CacheDirectory     string   `yaml:"cache"`
	Pid                string   `yaml:"pid"`
	Address            string   `yaml:"address"`
	Port               int      `yaml:"port"`
	Debug              bool     `yaml:"debug"`
}

func LoadConfig(path string) (config *Config, err error) {
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

func (self *Config) Save(path string) error {
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
func (self *Config) InitializeConfig(path string) error {
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
