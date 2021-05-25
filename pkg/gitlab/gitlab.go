package gitlab

import (
	"github.com/Spazzy757/credentials-rotator/pkg/config"
	"github.com/Spazzy757/credentials-rotator/pkg/helpers"
	"github.com/xanzy/go-gitlab"
)

func UpdateVariable(
	client *gitlab.Client,
	cred *config.Credential,
	value string,
) error {
	opts := &gitlab.UpdateProjectVariableOptions{
		Value:        gitlab.String(value),
		VariableType: gitlab.VariableType("file"),
	}

	_, _, err := client.ProjectVariables.UpdateVariable(
		cred.ProjectID,
		cred.Variable,
		opts,
	)
	return err
}

func GetClient() (*gitlab.Client, error) {
	token := helpers.GetEnv("GITLAB_TOKEN", "")
	c, err := gitlab.NewClient(token)
	if err != nil {
		return &gitlab.Client{}, err
	}
	return c, nil
}
