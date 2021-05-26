package handlers

import (
	iam "cloud.google.com/go/iam/admin/apiv1"
	"github.com/Spazzy757/credentials-rotator/pkg/config"
	"github.com/Spazzy757/credentials-rotator/pkg/gitlab"
	"github.com/Spazzy757/credentials-rotator/pkg/google"
)

func ConfigHandler(cfg *config.Config) error {
	var err error
	for _, cred := range cfg.Credentials {
		switch cred.Type {
		case "gitlab":
			err = gitlabHandler(cfg, &cred, cfg.GoogleIAMClient)
		}
	}
	return err
}

func gitlabHandler(
	cfg *config.Config,
	cred *config.Credential,
	client *iam.IamClient,
) error {
	key, err := google.CreateKey(
		cfg.Ctx,
		cred.GoogleProjectID,
		cred.ServiceAccount,
		client,
	)
	if err != nil {
		return err
	}
	err = gitlab.UpdateVariable(cfg.GitlabClient, cred, string(key.PrivateKeyData))
	return err
}
