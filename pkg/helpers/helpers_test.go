package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("SET_ENV", "1")
	t.Run("Test Unset Environment Returns Default", func(t *testing.T) {
		environment := GetEnv("UNSET_ENV", "default")
		assert.Equal(t, environment, "default")
	})
	t.Run("Test Set Environment Returns Default", func(t *testing.T) {
		environment := GetEnv("SET_ENV", "2")
		assert.Equal(t, environment, "1")
	})
}
