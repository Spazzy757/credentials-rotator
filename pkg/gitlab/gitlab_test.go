package gitlab

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Spazzy757/credentials-rotator/pkg/config"
	"github.com/Spazzy757/credentials-rotator/pkg/helpers"
	"github.com/stretchr/testify/require"
)

func TestUpdateVariable(t *testing.T) {
	t.Run("Test Varibale Gets Updated", func(t *testing.T) {
		assertions := require.New(t)
		mux, server, client := helpers.SetupGitlabTestServer(t)
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
		mux, server, client := helpers.SetupGitlabTestServer(t)
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
