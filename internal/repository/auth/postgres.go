package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateSession(ctx context.Context, session *auth.RefreshSession) error {
	query := `
		INSERT INTO refresh_sessions 
		(id, user_id, refresh_token, expires_at, issued_at, is_active, user_agent, ip_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	now := time.Now()
	session.CreatedAt = now
	session.UpdatedAt = now
	session.IsActive = true

	err := r.db.QueryRow(ctx, query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.ExpiresAt,
		session.IssuedAt,
		session.IsActive,
		session.UserAgent,
		session.IPAddress,
		session.CreatedAt,
		session.UpdatedAt,
	).Scan()

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetSessionByRefreshToken(ctx context.Context, token string) (*auth.RefreshSession, error) {
	query := `
		SELECT id, user_id, refresh_token, expires_at, issued_at, is_active, user_agent, ip_address, created_at, updated_at
		FROM refresh_sessions
		WHERE refresh_token = $1 AND is_active = true
	`

	var session auth.RefreshSession
	err := r.db.QueryRow(ctx, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAt,
		&session.IssuedAt,
		&session.IsActive,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("refresh session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	return &session, nil
}

func (r *PostgresRepository) GetSessionByID(ctx context.Context, id string) (*auth.RefreshSession, error) {
	query := `
		SELECT id, user_id, refresh_token, expires_at, issued_at, is_active, user_agent, ip_address, created_at, updated_at
		FROM refresh_sessions
		WHERE id = $1
	`

	var session auth.RefreshSession
	err := r.db.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAt,
		&session.IssuedAt,
		&session.IsActive,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("refresh session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (r *PostgresRepository) GetActiveSessionsByUserID(ctx context.Context, userID string) ([]*auth.RefreshSession, error) {
	query := `
		SELECT id, user_id, refresh_token, expires_at, issued_at, is_active, user_agent, ip_address, created_at, updated_at
		FROM refresh_sessions
		WHERE user_id = $1 AND is_active = true AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*auth.RefreshSession
	for rows.Next() {
		var session auth.RefreshSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.RefreshToken,
			&session.ExpiresAt,
			&session.IssuedAt,
			&session.IsActive,
			&session.UserAgent,
			&session.IPAddress,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return sessions, nil
}

func (r *PostgresRepository) UpdateSessionActivity(ctx context.Context, id string) error {
	query := `UPDATE refresh_sessions SET updated_at = NOW() WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (r *PostgresRepository) RevokeSession(ctx context.Context, id string) error {
	query := `UPDATE refresh_sessions SET is_active = false, updated_at = NOW() WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (r *PostgresRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	query := `UPDATE refresh_sessions SET is_active = false, updated_at = NOW() WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all sessions: %w", err)
	}
	return nil
}

func (r *PostgresRepository) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	query := `DELETE FROM refresh_sessions WHERE expires_at < NOW()`
	cmd, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	return cmd.RowsAffected(), nil
}
