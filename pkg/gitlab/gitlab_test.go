package gitlab

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/Spazzy757/credentials-rotator/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

func setup(t *testing.T) (*http.ServeMux, *httptest.Server, *gitlab.Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the Gitlab client being tested.
	client, err := gitlab.NewClient("", gitlab.WithBaseURL(server.URL))
	if err != nil {
		server.Close()
		t.Fatalf("Failed to create client: %v", err)
	}

	return mux, server, client
}

func TestGetClient(t *testing.T) {
	t.Run("Returns Client", func(t *testing.T) {
		client, _ := GetClient()
		expected := reflect.TypeOf(&gitlab.Client{})
		assert.Equal(t, reflect.TypeOf(client), expected)
	})
}

func TestUpdateVariable(t *testing.T) {
	t.Run("Test Varibale Gets Updated", func(t *testing.T) {
		assertions := require.New(t)
		mux, server, client := setup(t)
		defer server.Close()
		mux.HandleFunc("/api/v4/projects/12345/variables/TEST_VARIABLE",
			func(w http.ResponseWriter, r *http.Request) {
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
		os.Setenv("GITLAB_TOKEN", "XjZV85VTkeQvqwLEc8gb")
		creds := config.Credential{
			ProjectID: "12345",
			Variable:  "TEST_VARIABLE",
		}

		err := UpdateVariable(client, &creds, "ABCDBC")
		assertions.NoError(err)
	})
	t.Run("Test Varibale Fails", func(t *testing.T) {
		assertions := require.New(t)
		mux, server, client := setup(t)
		defer server.Close()
		mux.HandleFunc("/api/v4/projects/12345/variables/TEST_VARIABLE",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
		)
		os.Setenv("GITLAB_TOKEN", "XjZV85VTkeQvqwLEc8gb")
		creds := config.Credential{
			ProjectID: "12345",
			Variable:  "TEST_VARIABLE",
		}

		err := UpdateVariable(client, &creds, "ABCDBC")
		assertions.Error(err)
	})
}
