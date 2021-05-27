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

//Config used for the CLI command
type Config struct {
	Ctx context.Context

	// The GitlabClient that will be used to communicate
	// with a Gitlab instance
	GitlabClient *gitlab.Client

	// The Google IAM Client that is used to communicate with
	// Google Clouds IAM service
	GoogleIAMClient *iam.IamClient

	// List of credentials that will be used to update
	Credentials []Credential `yaml:"credentials,omitempty"`
}

//Credential that needs to be updated
type Credential struct {
	// The type of Credential (only supports gitlab for now)
	Type string `yaml:"type"`

	// The variable in the CI/CD to update
	// e.g GOOGLE_APPLICATION_CREDENTIAL
	Variable string `yaml:"variable"`

	// The Google Service Account email to update the key on
	ServiceAccount string `yaml:"service_account"`

	// Project ID the gitlab repos project ID
	ProjectID string `yaml:"project_id"`

	// Google Project ID where the service account is located
	GoogleProjectID string `yaml:"google_project_id"`
}

//LoadConfig loads the config from a file
//additionally adds clients to the config
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

//getGitlabClient creates a gitlab client
//and attaches it to the configuration
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
	// this should point to a test server
	serverURL := helpers.GetEnv("GITLAB_TEST_SERVER_URL", "")
	c, err := gitlab.NewClient("", gitlab.WithBaseURL(serverURL))
	if err != nil {
		return err
	}
	cfg.GitlabClient = c
	return nil
}

//getGoogleIAMClient creates a client and
//attaches it to the configuration
func getGoogleIAMClient(cfg *Config) error {
	c, err := iam.NewIamClient(cfg.Ctx)
	if err != nil {
		return err
	}
	cfg.GoogleIAMClient = c
	return nil
}
