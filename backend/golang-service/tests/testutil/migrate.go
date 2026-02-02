package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// RunMigrations runs database migrations using golang-migrate.
// It expects the migrate binary to be available or uses the migrations directory.
func RunMigrations(t *testing.T, databaseURL string) {
	t.Helper()

	migrationsDir := getMigrationsDir()

	// Try using golang-migrate CLI
	cmd := exec.Command("migrate",
		"-path", migrationsDir,
		"-database", databaseURL,
		"up",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Log output for debugging
		t.Logf("migrate output: %s", string(output))
		t.Fatalf("failed to run migrations: %v", err)
	}
}

// ResetDatabase drops and recreates all tables by running down then up migrations.
func ResetDatabase(t *testing.T, databaseURL string) {
	t.Helper()

	migrationsDir := getMigrationsDir()

	// Run down migrations (ignore errors if tables don't exist)
	downCmd := exec.Command("migrate",
		"-path", migrationsDir,
		"-database", databaseURL,
		"down", "-all",
	)
	_, _ = downCmd.CombinedOutput()

	// Run up migrations
	upCmd := exec.Command("migrate",
		"-path", migrationsDir,
		"-database", databaseURL,
		"up",
	)

	output, err := upCmd.CombinedOutput()
	if err != nil {
		t.Logf("migrate up output: %s", string(output))
		t.Fatalf("failed to run up migrations: %v", err)
	}
}

// getMigrationsDir returns the path to the migrations directory.
func getMigrationsDir() string {
	// Check if MIGRATIONS_DIR env var is set
	if dir := os.Getenv("MIGRATIONS_DIR"); dir != "" {
		return dir
	}

	// Default relative path (from tests directory or project root)
	paths := []string{
		"../../migrations",
		"../migrations",
		"./migrations",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return "migrations"
}

// MustRunMigrationsWithPool runs migrations using pgx directly without CLI.
// This is useful when the migrate CLI is not available.
func MustRunMigrationsWithPool(t *testing.T, databaseURL string) {
	t.Helper()

	// Read migration file
	migrationsDir := getMigrationsDir()
	upFile := fmt.Sprintf("%s/000001_init.up.sql", migrationsDir)

	content, err := os.ReadFile(upFile)
	if err != nil {
		t.Fatalf("failed to read migration file: %v", err)
	}

	// Execute using psql or directly
	cmd := exec.Command("psql", databaseURL, "-c", string(content))
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("psql output: %s", string(output))
		t.Fatalf("failed to run migration: %v", err)
	}
}
