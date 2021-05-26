package handlers

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"testing"

	iam "cloud.google.com/go/iam/admin/apiv1"
	"github.com/Spazzy757/credentials-rotator/pkg/config"
	"github.com/Spazzy757/credentials-rotator/pkg/helpers"
	"github.com/Spazzy757/credentials-rotator/pkg/test"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

var clientOpt option.ClientOption

func GetTestConfig(creds []config.Credential) config.Config {
	testConfig := config.Config{
		Credentials: creds,
	}
	configBytes, _ := yaml.Marshal(testConfig)
	tmpDir := os.TempDir()
	ioutil.WriteFile(path.Join(tmpDir, "config.yaml"), configBytes, 0644)
	defer os.RemoveAll(path.Join(tmpDir, "config.yaml"))
	cfg := config.Config{}
	_ = cfg.LoadConfig(path.Join(tmpDir, "config.yaml"))
	return cfg
}

var (
	mockIam test.MockIamServer
)

func mockGRPCServer() *grpc.Server {
	flag.Parse()

	serv := grpc.NewServer()
	adminpb.RegisterIAMServer(serv, &mockIam)

	lis, err := net.Listen("tcp", "localhost:5763")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := serv.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Printf("HERE: %v", err.Error())
	}
	clientOpt = option.WithGRPCConn(conn)

	return serv
}

func TestGitlabHandler(t *testing.T) {
	t.Run("handles gitlab", func(t *testing.T) {
		assertions := require.New(t)
		grpcServ := mockGRPCServer()
		defer grpcServ.GracefulStop()
		key := adminpb.ServiceAccountKey{}
		key.PrivateKeyData = []byte(`
    {
        "type": "service_account",
	      "project_id": "test-0000000",
		    "private_key_id": "00000000000000000000000000000000000000000000000000",
			  "private_key": "-----BEGIN PRIVATE KEY-----\n\n-----END PRIVATE KEY-----\n",
				"client_email": "test@test-0000000.iam.gserviceaccount.com",
				"client_id": "000000000000000000000",
				"auth_uri": "https://accounts.google.com/o/oauth2/auth",
				"token_uri": "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/test%40test-00000.iam.gserviceaccount.com"
		}`)
		expectedResponse := &key
		mockIam.Err = nil
		mockIam.Reqs = nil
		mockIam.Resps = append(mockIam.Resps[:0], expectedResponse)

		os.Setenv("TEST", "true")
		mux, server, _ := helpers.SetupGitlabTestServer(t)
		mux.HandleFunc("/api/v4/projects/12345/variables/TEST_VARIABLE",
			func(w http.ResponseWriter, r *http.Request) {
				_, err := ioutil.ReadAll(r.Body)
				assertions.NoError(err)
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{
					    "key": "TEST_VARIABLE",
							"value": "updated value",
							"variable_type": "file",
							"protected": true,
							"masked": false,
							"environment_scope": "*"
				}`)
			},
		)
		defer server.Close()
		os.Setenv("GITLAB_TEST_SERVER_URL", server.URL)
		creds := []config.Credential{
			config.Credential{
				Type:            "gitlab",
				ProjectID:       "12345",
				Variable:        "TEST_VARIABLE",
				GoogleProjectID: "test-0000000",
				ServiceAccount:  "test@test-0000000.iam.gserviceaccount.com",
			},
		}
		c, _ := iam.NewIamClient(context.Background(), clientOpt)
		cfg := GetTestConfig(creds)
		err := gitlabHandler(&cfg, &cfg.Credentials[0], c)
		assertions.NoError(err)
	})
}
