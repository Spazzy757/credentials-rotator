package gitlab

import (
	"github.com/Spazzy757/credentials-rotator/pkg/config"
	"github.com/xanzy/go-gitlab"
)

//UpdateVariable takes a string and updates the
//specified CI/CD variable that is in the
//cred struct
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
