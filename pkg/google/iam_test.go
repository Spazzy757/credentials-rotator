package google

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"

	iam "cloud.google.com/go/iam/admin/apiv1"
	"github.com/golang/protobuf/ptypes"
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
	return &adminpb.ServiceAccountKey{}, nil
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

		client, err := iam.NewIamClient(ctx, clientOpt)

		assertions.NoError(err)

		key, err := CreateKey(ctx, project, serviceAccount, client)

		expectedKey := reflect.TypeOf(&adminpb.ServiceAccountKey{})
		assertions.NoError(err)
		assertions.Equal(expectedKey, reflect.TypeOf(key))
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
