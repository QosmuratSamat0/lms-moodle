//go:build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/ap1-final-mini-moodle/tests/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

// TestMain sets up and tears down the test database.
func TestMain(m *testing.M) {
	// Get test config
	cfg := testutil.DefaultTestConfig()

	// Connect to test database with retry
	ctx := context.Background()
	var err error
	testPool, err = pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	defer testPool.Close()

	// Run migrations (if using migrate CLI)
	// Alternatively, you can run migrations via docker-compose before tests
	// testutil.RunMigrations(nil, cfg.DatabaseURL) // Note: can't use t here

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// getTestPool returns the shared test database pool.
func getTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testPool == nil {
		t.Fatal("test database pool not initialized")
	}
	return testPool
}
