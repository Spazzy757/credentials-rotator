package google

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	iam "cloud.google.com/go/iam/admin/apiv1"
	"github.com/golang/protobuf/ptypes"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var _ = io.EOF
var _ = ptypes.MarshalAny
var _ status.Status

type mockIamServer struct {
	// Embed for forward compatibility.
	// Tests will keep working if more methods are added
	// in the future.
	adminpb.IAMServer

	reqs []proto.Message

	// If set, all calls return this error.
	err error

	// responses to return if err == nil
	resps []proto.Message
}

func (s *mockIamServer) CreateServiceAccountKey(
	ctx context.Context,
	req *adminpb.CreateServiceAccountKeyRequest,
) (*adminpb.ServiceAccountKey, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*adminpb.ServiceAccountKey), nil
}

func (s *mockIamServer) ListServiceAccountKeys(
	ctx context.Context,
	req *adminpb.ListServiceAccountKeysRequest,
) (*adminpb.ListServiceAccountKeysResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*adminpb.ListServiceAccountKeysResponse), nil
}

func (s *mockIamServer) DeleteServiceAccountKey(
	ctx context.Context,
	req *adminpb.DeleteServiceAccountKeyRequest,
) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		return nil, s.err
	}
	return s.resps[0].(*emptypb.Empty), nil
}

// clientOpt is the option tests should use to connect to the test server.
// It is initialized by TestMain.
var clientOpt option.ClientOption

var (
	mockIam mockIamServer
)

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

func TestGetKeys(t *testing.T) {
	t.Run("Test Get Keys", func(t *testing.T) {
		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"
		expectedResponse := &adminpb.ServiceAccountKey{}
		mockIam.err = nil
		mockIam.reqs = nil
		mockIam.resps = append(mockIam.resps[:0], expectedResponse)

		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)

		key, err := CreateKey(ctx, project, serviceAccount, client)

		assertions.NoError(err)
		assertions.Equal(expectedResponse, key)
	})
	t.Run("Test Get Keys Returned Error", func(t *testing.T) {
		errCode := codes.PermissionDenied
		mockIam.err = gstatus.Error(errCode, "test error")

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
	t.Run("Test List Keys", func(t *testing.T) {
		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"

		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)
		expectedResponse := &adminpb.ListServiceAccountKeysResponse{}
		mockIam.err = nil
		mockIam.reqs = nil
		mockIam.resps = append(mockIam.resps[:0], expectedResponse)

		key, err := ListKeys(ctx, project, serviceAccount, client)

		assertions.NoError(err)
		assertions.Equal(expectedResponse, key)
	})
	t.Run("Test List Keys Returned Error", func(t *testing.T) {
		errCode := codes.PermissionDenied
		mockIam.err = gstatus.Error(errCode, "test error")

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
	t.Run("Test Delete Key", func(t *testing.T) {
		assertions := require.New(t)
		ctx := context.Background()
		project := "project-1"
		serviceAccount := "test@example.com"
		key := "key-name"
		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)
		expectedResponse := &emptypb.Empty{}
		mockIam.err = nil
		mockIam.reqs = nil
		mockIam.resps = append(mockIam.resps[:0], expectedResponse)

		err = DeleteKey(ctx, project, serviceAccount, key, client)

		assertions.NoError(err)
	})
	t.Run("Test Delete Key Returned Error", func(t *testing.T) {
		errCode := codes.PermissionDenied
		mockIam.err = gstatus.Error(errCode, "test error")

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
