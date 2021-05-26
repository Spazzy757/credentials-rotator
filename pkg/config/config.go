package config

import (
	"context"
	"fmt"
	"io/ioutil"

	iam "cloud.google.com/go/iam/admin/apiv1"
	"github.com/Spazzy757/credentials-rotator/pkg/helpers"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GitlabClient    *gitlab.Client
	GoogleIAMClient *iam.IamClient
	Ctx             context.Context
	Credentials     []Credential `yaml:"credentials,omitempty"`
}

type Credential struct {
	Type            string `yaml:"type"`
	Variable        string `yaml:"variable"`
	ServiceAccount  string `yaml:"service_account"`
	ProjectID       string `yaml:"project_id"`
	GoogleProjectID string `yaml:"google_project_id"`
}

//LoadConfig loads the config from a file
func (c *Config) LoadConfig(config_file string) error {
	config, err := ioutil.ReadFile(config_file)
	if err != nil {
		return fmt.Errorf("failed loading configuration")
	}
	err = yaml.Unmarshal(config, c)
	if err != nil {
		return err
	}
	ctx := context.Background()
	c.Ctx = ctx
	err = getGitlabClient(c)
	if err != nil {
		return err
	}
	err = getGoogleIAMClient(c)
	return err
}

func getGitlabClient(cfg *Config) error {
	isTest := helpers.GetEnv("TEST", "")
	if isTest != "true" {
		token := helpers.GetEnv("GITLAB_TOKEN", "")
		c, err := gitlab.NewClient(token)
		if err != nil {
			return err
		}
		cfg.GitlabClient = c
		return nil
	}
	// If test check for test URL
	serverURL := helpers.GetEnv("GITLAB_TEST_SERVER_URL", "")
	client, err := gitlab.NewClient("", gitlab.WithBaseURL(serverURL))
	if err != nil {
		return err
	}
	cfg.GitlabClient = client
	return nil
}

func getGoogleIAMClient(cfg *Config) error {
	c, err := iam.NewIamClient(cfg.Ctx)
	if err != nil {
		return err
	}
	cfg.GoogleIAMClient = c
	return nil
}
