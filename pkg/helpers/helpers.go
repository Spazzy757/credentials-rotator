package helpers

import (
	"os"
)

//GetEnv gets an environment variable or returns default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
