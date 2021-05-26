package test

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/golang/protobuf/ptypes"
	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ = io.EOF
var _ = ptypes.MarshalAny
var _ status.Status

type MockIamServer struct {
	// Embed for forward compatibility.
	// Tests will keep working if more methods are added
	// in the future.
	adminpb.IAMServer

	Reqs []proto.Message

	// If set, all calls return this error.
	Err error

	// responses to return if err == nil
	Resps []proto.Message
}

func (s *MockIamServer) CreateServiceAccountKey(
	ctx context.Context,
	req *adminpb.CreateServiceAccountKeyRequest,
) (*adminpb.ServiceAccountKey, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.Reqs = append(s.Reqs, req)
	if s.Err != nil {
		return nil, s.Err
	}
	return s.Resps[0].(*adminpb.ServiceAccountKey), nil
}

func (s *MockIamServer) ListServiceAccountKeys(
	ctx context.Context,
	req *adminpb.ListServiceAccountKeysRequest,
) (*adminpb.ListServiceAccountKeysResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.Reqs = append(s.Reqs, req)
	if s.Err != nil {
		return nil, s.Err
	}
	return s.Resps[0].(*adminpb.ListServiceAccountKeysResponse), nil
}

func (s *MockIamServer) DeleteServiceAccountKey(
	ctx context.Context,
	req *adminpb.DeleteServiceAccountKeyRequest,
) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if xg := md["x-goog-api-client"]; len(xg) == 0 || !strings.Contains(xg[0], "gl-go/") {
		return nil, fmt.Errorf("x-goog-api-client = %v, expected gl-go key", xg)
	}
	s.Reqs = append(s.Reqs, req)
	if s.Err != nil {
		return nil, s.Err
	}
	return s.Resps[0].(*emptypb.Empty), nil
}
