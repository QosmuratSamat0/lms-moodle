//go:build e2e

package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/ap1-final-mini-moodle/tests/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testPool *pgxpool.Pool
	baseURL  string
)

// TestMain sets up the e2e test environment.
func TestMain(m *testing.M) {
	cfg := testutil.DefaultTestConfig()

	// Get API base URL from environment
	baseURL = os.Getenv("TEST_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Connect to test database
	ctx := context.Background()
	var err error
	testPool, err = pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	defer testPool.Close()

	code := m.Run()
	os.Exit(code)
}

func getTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testPool == nil {
		t.Fatal("test database pool not initialized")
	}
	return testPool
}

func getBaseURL() string {
	return baseURL
}
