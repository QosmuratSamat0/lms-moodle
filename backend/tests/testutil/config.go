package testutil

import (
	"os"

	"github.com/ap1-final-mini-moodle/internal/shared/config"
)

// GetTestConfig returns test database configuration
func GetTestConfig() *config.Config {
	// Set test environment variables
	os.Setenv("ENV", "test")
	os.Setenv("DB_HOST", os.Getenv("TEST_DB_HOST"))
	if os.Getenv("TEST_DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}

	os.Setenv("DB_PORT", os.Getenv("TEST_DB_PORT"))
	if os.Getenv("TEST_DB_PORT") == "" {
		os.Setenv("DB_PORT", "5432")
	}

	os.Setenv("DB_NAME", os.Getenv("TEST_DB_NAME"))
	if os.Getenv("TEST_DB_NAME") == "" {
		os.Setenv("DB_NAME", "moodle_test")
	}

	os.Setenv("DB_USER", os.Getenv("TEST_DB_USER"))
	if os.Getenv("TEST_DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
	}

	os.Setenv("DB_PASSWORD", os.Getenv("TEST_DB_PASSWORD"))
	if os.Getenv("TEST_DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "postgres")
	}

	return config.Load()
}
