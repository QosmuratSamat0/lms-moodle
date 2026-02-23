package auth

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/auth"
	"github.com/ap1-final-mini-moodle/internal/shared/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db    *pgxpool.Pool
	redis *database.RedisClient
}

func NewPostgresRepository(db *pgxpool.Pool, redis *database.RedisClient) *PostgresRepository {
	return &PostgresRepository{db: db, redis: redis}
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

	_, err := r.db.Exec(ctx, query,
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
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	hashed := sha256.Sum256([]byte(session.RefreshToken))
	key := fmt.Sprintf("session:refresh:%x", hashed)
	data, err := json.Marshal(session)
	expiration := time.Until(session.ExpiresAt)

	if r.redis != nil {
		_ = r.redis.Set(ctx, "session:"+session.ID, data, expiration)
		_ = r.redis.Set(ctx, key, data, expiration)
		_ = r.redis.Del(ctx, "sessions:user:"+session.UserID)
	}

	return nil
}

func (r *PostgresRepository) GetSessionByRefreshToken(ctx context.Context, token string) (*auth.RefreshSession, error) {
	ctx = context.Background()
	cacheKey := fmt.Sprintf("session:refresh:%x", sha256.Sum256([]byte(token)))

	query := `
		SELECT id, user_id, refresh_token, expires_at, issued_at, is_active, user_agent, ip_address, created_at, updated_at
		FROM refresh_sessions
		WHERE refresh_token = $1 AND is_active = true
	`

	if r.redis != nil {
		cached, err := r.redis.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var session auth.RefreshSession
			if err := json.Unmarshal([]byte(cached), &session); err == nil {
				if time.Now().After(session.ExpiresAt) {
					return nil, fmt.Errorf("refresh token expired")
				}
				return &session, nil
			}
		}
	}

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

	data, err := json.Marshal(session)
	expiration := time.Until(session.ExpiresAt)
	if err == nil {
		_ = r.redis.Set(ctx, cacheKey, data, expiration)
	}

	return &session, nil
}

func (r *PostgresRepository) GetSessionByID(ctx context.Context, id string) (*auth.RefreshSession, error) {
	ctx = context.Background()
	cacheKey := "session:" + id
	cached, err := r.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var session auth.RefreshSession
		if err := json.Unmarshal([]byte(cached), &session); err == nil {
			if time.Now().After(session.ExpiresAt) {
				return nil, fmt.Errorf("refresh token expired")
			}
			return &session, nil
		}
	}
	query := `
		SELECT id, user_id, refresh_token, expires_at, issued_at, is_active, user_agent, ip_address, created_at, updated_at
		FROM refresh_sessions
		WHERE id = $1 AND is_active = true

	`

	var session auth.RefreshSession
	err = r.db.QueryRow(ctx, query, id).Scan(
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

	data, err := json.Marshal(session)
	expiration := time.Until(session.ExpiresAt)
	if err == nil {
		_ = r.redis.Set(ctx, cacheKey, data, expiration)
	}

	return &session, nil
}

func (r *PostgresRepository) GetActiveSessionsByUserID(ctx context.Context, userID string) ([]*auth.RefreshSession, error) {
	ctx = context.Background()
	cacheKey := "sessions:user:" + userID
	cached, err := r.redis.Get(ctx, cacheKey)
	query := `
		SELECT id, user_id, refresh_token, expires_at, issued_at, is_active, user_agent, ip_address, created_at, updated_at
		FROM refresh_sessions
		WHERE user_id = $1 AND is_active = true AND expires_at > NOW()
		ORDER BY created_at DESC
	`
	if err == nil && cached != "" {
		var sessions []*auth.RefreshSession
		if err := json.Unmarshal([]byte(cached), &sessions); err == nil {
			return sessions, nil
		}
	}

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

	data, err := json.Marshal(sessions)
	if err == nil {
		_ = r.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return sessions, nil
}

func (r *PostgresRepository) UpdateSessionActivity(ctx context.Context, id string) error {

	var refreshToken string
	var userID string
	err := r.db.QueryRow(ctx,
		`SELECT refresh_token, user_id FROM refresh_sessions WHERE id = $1`,
		id,
	).Scan(&refreshToken, &userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("session not found")
		}
		return err
	}

	cmd, err := r.db.Exec(ctx,
		`UPDATE refresh_sessions SET updated_at = NOW() WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	cacheKeyID := "session:" + id
	cacheKeyRefresh := fmt.Sprintf(
		"session:refresh:%x",
		sha256.Sum256([]byte(refreshToken)),
	)

	_ = r.redis.Del(ctx, cacheKeyID)
	_ = r.redis.Del(ctx, cacheKeyRefresh)
	_ = r.redis.Del(ctx, "sessions:user:"+userID)

	return nil
}

func (r *PostgresRepository) RevokeSession(ctx context.Context, id string) error {
	var refreshToken string
	var userID string
	err := r.db.QueryRow(ctx,
		`SELECT refresh_token, user_id FROM refresh_sessions WHERE id = $1`,
		id,
	).Scan(&refreshToken, &userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("session not found")
		}
		return err
	}

	cmd, err := r.db.Exec(ctx,
		`UPDATE refresh_sessions 
		 SET is_active = false, updated_at = NOW() 
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	cacheKeyID := "session:" + id
	cacheKeyRefresh := fmt.Sprintf(
		"session:refresh:%x",
		sha256.Sum256([]byte(refreshToken)),
	)

	_ = r.redis.Del(ctx, cacheKeyID)
	_ = r.redis.Del(ctx, cacheKeyRefresh)
	_ = r.redis.Del(ctx, "sessions:user:"+userID)

	return nil
}

func (r *PostgresRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {

	cmd, err := r.db.Exec(ctx,
		`UPDATE refresh_sessions
		 SET is_active = false, updated_at = NOW()
		 WHERE user_id = $1`,
		userID,
	)
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}
	if err != nil {
		return fmt.Errorf("failed to revoke all sessions: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeleteExpiredSessions(ctx context.Context) (int64, error) {

	rows, err := r.db.Query(ctx,
		`SELECT id, refresh_token 
		 FROM refresh_sessions 
		 WHERE expires_at < NOW()`,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	type sessionData struct {
		ID           string
		RefreshToken string
	}

	var sessions []sessionData

	for rows.Next() {
		var s sessionData
		if err := rows.Scan(&s.ID, &s.RefreshToken); err != nil {
			return 0, err
		}
		sessions = append(sessions, s)
	}

	cmd, err := r.db.Exec(ctx,
		`DELETE FROM refresh_sessions WHERE expires_at < NOW()`,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	for _, s := range sessions {
		cacheKeyID := "session:" + s.ID
		cacheKeyRefresh := fmt.Sprintf(
			"session:refresh:%x",
			sha256.Sum256([]byte(s.RefreshToken)),
		)

		_ = r.redis.Del(ctx, cacheKeyID)
		_ = r.redis.Del(ctx, cacheKeyRefresh)
	}

	return cmd.RowsAffected(), nil
}
