package config

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestConfigLoad(t *testing.T) {
	t.Run("Testing Loading File Returns Configuration", func(t *testing.T) {
		assertions := require.New(t)
		testConfig := Config{
			Credentials: []Credential{
				Credential{
					Type:           "test",
					Project:        "foo/bar",
					Variable:       "TEST_VARIABLE",
					ServiceAccount: "test@example.com",
				},
			},
		}
		configBytes, err := yaml.Marshal(testConfig)
		assertions.NoError(err)
		tmpDir := os.TempDir()
		ioutil.WriteFile(path.Join(tmpDir, "config.yaml"), configBytes, 0644)
		defer os.RemoveAll(path.Join(tmpDir, "config.yaml"))
		cfg := Config{}
		err = cfg.LoadConfig(path.Join(tmpDir, "config.yaml"))
		assertions.NoError(err)
		assertions.Equal(cfg.Credentials[0].Type, "test")
		assertions.Equal(cfg.Credentials[0].Project, "foo/bar")
		assertions.Equal(cfg.Credentials[0].Variable, "TEST_VARIABLE")
		assertions.Equal(cfg.Credentials[0].ServiceAccount, "test@example.com")
	})
	t.Run("Testing Loading Invalid File", func(t *testing.T) {
		assertions := require.New(t)
		tmpDir := os.TempDir()
		configBytes := []byte("`^88(0")
		ioutil.WriteFile(path.Join(tmpDir, "config.yaml"), configBytes, 0644)
		defer os.RemoveAll(path.Join(tmpDir, "config.yaml"))
		cfg := Config{}
		err := cfg.LoadConfig(path.Join(tmpDir, "config.yaml"))
		assertions.Error(err)
		assertions.Equal(cfg.Credentials, []Credential(nil))
	})
	t.Run("Testing Loading Non Existant File", func(t *testing.T) {
		assertions := require.New(t)
		cfg := Config{}
		err := cfg.LoadConfig("config.yaml")
		assertions.Error(err)
		assertions.Equal(cfg.Credentials, []Credential(nil))
	})
}
