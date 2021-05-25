package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Credentials []Credential `yaml:credentials,omitempty`
}

type Credential struct {
	Type           string `yaml:type`
	Project        string `yaml:project`
	Variable       string `yaml:variable`
	ServiceAccount string `yaml:service_account`
	ProjectID      string `yaml:project_id`
}

//LoadConfig loads the config from a file
func (c *Config) LoadConfig(config_file string) error {
	config, err := ioutil.ReadFile(config_file)
	if err != nil {
		return fmt.Errorf("failed loading configuration")
	}
	err = yaml.Unmarshal(config, c)
	return err
}
