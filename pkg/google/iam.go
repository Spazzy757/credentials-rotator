package google

import (
	"context"
	"fmt"

	iam "cloud.google.com/go/iam/admin/apiv1"
	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"
)

func CreateKey(
	ctx context.Context,
	project string,
	serviceAccount string,
	client *iam.IamClient,
) (*adminpb.ServiceAccountKey, error) {
	formatted := fmt.Sprintf(
		"projects/%v/serviceAccounts/%v",
		project,
		serviceAccount,
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
