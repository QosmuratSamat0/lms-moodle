package testutil

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxTest wraps a test function with a transaction that is rolled back after the test.
// This provides test isolation without needing to truncate tables.
type TxTest struct {
	pool *pgxpool.Pool
}

// NewTxTest creates a new TxTest helper.
func NewTxTest(pool *pgxpool.Pool) *TxTest {
	return &TxTest{pool: pool}
}

// Run executes the test function within a transaction that is rolled back after completion.
func (tt *TxTest) Run(t *testing.T, name string, fn func(t *testing.T, tx pgx.Tx)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		ctx := context.Background()
		tx, err := tt.pool.Begin(ctx)
		if err != nil {
			t.Fatalf("failed to begin transaction: %v", err)
		}

		defer func() {
			if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
				t.Logf("warning: failed to rollback transaction: %v", err)
			}
		}()

		fn(t, tx)
	})
}

// RunWithContext executes the test function within a transaction with a provided context.
func (tt *TxTest) RunWithContext(t *testing.T, name string, ctx context.Context, fn func(t *testing.T, ctx context.Context, tx pgx.Tx)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		tx, err := tt.pool.Begin(ctx)
		if err != nil {
			t.Fatalf("failed to begin transaction: %v", err)
		}

		defer func() {
			if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
				t.Logf("warning: failed to rollback transaction: %v", err)
			}
		}()

		fn(t, ctx, tx)
	})
}
