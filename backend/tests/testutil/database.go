package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MustOpenTestDB opens a connection to the test database.
// It fails the test if the connection cannot be established.
func MustOpenTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg := DefaultTestConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("failed to ping test database: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

// MustOpenTestDBWithRetry opens a connection to the test database with retry logic.
// Useful when waiting for Docker container to be ready.
func MustOpenTestDBWithRetry(t *testing.T, maxRetries int, retryDelay time.Duration) *pgxpool.Pool {
	t.Helper()

	cfg := DefaultTestConfig()
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err = pgxpool.New(ctx, cfg.DatabaseURL)
		cancel()

		if err == nil {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
			err = pool.Ping(pingCtx)
			pingCancel()

			if err == nil {
				t.Cleanup(func() {
					pool.Close()
				})
				return pool
			}
			pool.Close()
		}

		if i < maxRetries-1 {
			t.Logf("Retry %d/%d: waiting for database... (%v)", i+1, maxRetries, err)
			time.Sleep(retryDelay)
		}
	}

	t.Fatalf("failed to connect to test database after %d retries: %v", maxRetries, err)
	return nil
}

// TruncateAll truncates all tables in the test database.
// Use this to reset state between tests when not using per-test transactions.
func TruncateAll(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tables := []string{
		"analytics_events",
		"schedule_events",
		"notifications",
		"chat_messages",
		"chat_room_members",
		"chat_rooms",
		"attendance_marks",
		"attendance_sessions",
		"grades",
		"submissions",
		"assignments",
		"enrollments",
		"courses",
		"managers",
		"teachers",
		"students",
		"users",
	}

	for _, table := range tables {
		_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}

// TruncateTables truncates specific tables.
func TruncateTables(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, table := range tables {
		_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}
