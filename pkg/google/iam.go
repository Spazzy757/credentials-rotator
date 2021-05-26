package google

import (
	"context"
	"fmt"
	"net/url"

	iam "cloud.google.com/go/iam/admin/apiv1"
	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"
)

//CreateKey creates a new key for a service account
func CreateKey(
	ctx context.Context,
	project string,
	serviceAccount string,
	client *iam.IamClient,
) (*adminpb.ServiceAccountKey, error) {
	formatted := fmt.Sprintf(
		"projects/%v/serviceAccounts/%v",
		project,
		url.QueryEscape(serviceAccount),
	)
	request := &adminpb.CreateServiceAccountKeyRequest{
		Name: formatted,
	}
	resp, err := client.CreateServiceAccountKey(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//ListKeys returns all the keys under a service account
func ListKeys(
	ctx context.Context,
	project string,
	serviceAccount string,
	client *iam.IamClient,
) (*adminpb.ListServiceAccountKeysResponse, error) {
	formatted := fmt.Sprintf(
		"projects/%v/serviceAccounts/%v",
		project,
		serviceAccount,
	)
	request := &adminpb.ListServiceAccountKeysRequest{
		Name: formatted,
	}
	resp, err := client.ListServiceAccountKeys(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//DeleteKey deletes a service account key
func DeleteKey(
	ctx context.Context,
	project string,
	serviceAccount string,
	key string,
	client *iam.IamClient,
) error {
	formatted := fmt.Sprintf(
		"projects/%s/serviceAccounts/%s/keys/%s",
		project,
		serviceAccount,
		key,
	)
	request := &adminpb.DeleteServiceAccountKeyRequest{
		Name: formatted,
	}
	err := client.DeleteServiceAccountKey(ctx, request)
	return err
}
