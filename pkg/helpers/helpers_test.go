package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("SET_ENV", "1")
	t.Run("unset environment returns default", func(t *testing.T) {
		environment := GetEnv("UNSET_ENV", "default")
		assert.Equal(t, environment, "default")
	})
	t.Run("set environment returns environment", func(t *testing.T) {
		environment := GetEnv("SET_ENV", "2")
		assert.Equal(t, environment, "1")
	})
}
