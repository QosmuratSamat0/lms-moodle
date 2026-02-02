// Package testutil provides testing utilities for the aLMS project.
package testutil

import (
	"os"
	"strconv"
	"time"
)

// TestConfig holds configuration for test environment.
type TestConfig struct {
	DatabaseURL string
	RedisURL    string
	Timeout     time.Duration
}

// DefaultTestConfig returns the default test configuration
// with values from environment variables or sensible defaults.
func DefaultTestConfig() TestConfig {
	return TestConfig{
		DatabaseURL: getEnv("TEST_DATABASE_URL", "postgres://alms_test:alms_test_secret@localhost:5433/alms_test_db?sslmode=disable"),
		RedisURL:    getEnv("TEST_REDIS_URL", "redis://localhost:6380/0"),
		Timeout:     getDurationEnv("TEST_TIMEOUT", 30*time.Second),
	}
}

// getEnv returns the value of the environment variable or the default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv returns the duration value of the environment variable or the default value.
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}
