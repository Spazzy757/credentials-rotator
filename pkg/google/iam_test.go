package google

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"testing"

	iam "cloud.google.com/go/iam/admin/apiv1"
	"github.com/Spazzy757/credentials-rotator/pkg/test"
	"github.com/golang/protobuf/ptypes"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

var _ = io.EOF
var _ = ptypes.MarshalAny
var _ status.Status

// clientOpt is the option tests should use to connect to the test server.
// It is initialized by TestMain.
var clientOpt option.ClientOption

var (
	mockIam test.MockIamServer
)

//TestMain Setups the mock GRPC Server
func TestMain(m *testing.M) {
	flag.Parse()

	serv := grpc.NewServer()
	adminpb.RegisterIAMServer(serv, &mockIam)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	go serv.Serve(lis)

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	clientOpt = option.WithGRPCConn(conn)

	os.Exit(m.Run())
}

func TestCreateKey(t *testing.T) {
	t.Run("create key is successful", func(t *testing.T) {
		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"
		expectedResponse := &adminpb.ServiceAccountKey{}
		mockIam.Err = nil
		mockIam.Reqs = nil
		mockIam.Resps = append(mockIam.Resps[:0], expectedResponse)

		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)

		key, err := CreateKey(ctx, project, serviceAccount, client)

		assertions.NoError(err)
		assertions.Equal(expectedResponse, key)
	})
	t.Run("create key fails with error", func(t *testing.T) {
		errCode := codes.PermissionDenied
		mockIam.Err = gstatus.Error(errCode, "test error")

		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"

		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)

		_, err = CreateKey(ctx, project, serviceAccount, client)
		assertions.Error(err)
	})
}

func TestListKeys(t *testing.T) {
	t.Run("list keys is successful", func(t *testing.T) {
		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"

		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)
		expectedResponse := &adminpb.ListServiceAccountKeysResponse{}
		mockIam.Err = nil
		mockIam.Reqs = nil
		mockIam.Resps = append(mockIam.Resps[:0], expectedResponse)

		key, err := ListKeys(ctx, project, serviceAccount, client)

		assertions.NoError(err)
		assertions.Equal(expectedResponse, key)
	})
	t.Run("list keys fails with error", func(t *testing.T) {
		errCode := codes.PermissionDenied
		mockIam.Err = gstatus.Error(errCode, "test error")

		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"

		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)

		_, err = ListKeys(ctx, project, serviceAccount, client)
		assertions.Error(err)
	})
}

func TestDeleteKey(t *testing.T) {
	t.Run("delete key is successful", func(t *testing.T) {
		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"
		key := "key-name"
		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)
		expectedResponse := &emptypb.Empty{}
		mockIam.Err = nil
		mockIam.Reqs = nil
		mockIam.Resps = append(mockIam.Resps[:0], expectedResponse)

		err = DeleteKey(ctx, project, serviceAccount, key, client)

		assertions.NoError(err)
	})
	t.Run("delete key fails with error", func(t *testing.T) {
		errCode := codes.PermissionDenied
		mockIam.Err = gstatus.Error(errCode, "test error")

		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"
		key := "key-name"
		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)

		err = DeleteKey(ctx, project, serviceAccount, key, client)
		assertions.Error(err)
	})
}
