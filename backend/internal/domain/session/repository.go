package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the session repository interface
type Repository interface {
	Create(ctx context.Context, session *Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (*SessionWithUser, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Session, error)
	RevokeByID(ctx context.Context, id uuid.UUID) error
	RevokeByRefreshToken(ctx context.Context, refreshToken string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Session, int64, error)
	DeleteExpired(ctx context.Context) (int64, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new session repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		session.UserID, session.RefreshToken, session.UserAgent, session.IPAddress, session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *repository) GetByRefreshToken(ctx context.Context, refreshToken string) (*SessionWithUser, error) {
	query := `
		SELECT 
			s.id, s.user_id, s.refresh_token, s.user_agent, s.ip_address, s.expires_at, s.revoked_at, s.created_at,
			u.email AS user_email, u.role AS user_role, u.is_active AS user_is_active
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.refresh_token = $1`

	var sess SessionWithUser
	err := r.db.QueryRow(ctx, query, refreshToken).Scan(
		&sess.ID, &sess.UserID, &sess.RefreshToken, &sess.UserAgent, &sess.IPAddress, &sess.ExpiresAt, &sess.RevokedAt, &sess.CreatedAt,
		&sess.UserEmail, &sess.UserRole, &sess.UserIsActive,
	)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at, revoked_at, created_at
		FROM sessions
		WHERE id = $1`

	var sess Session
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sess.ID, &sess.UserID, &sess.RefreshToken, &sess.UserAgent, &sess.IPAddress, &sess.ExpiresAt, &sess.RevokedAt, &sess.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (r *repository) RevokeByID(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE sessions SET revoked_at = now() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) RevokeByRefreshToken(ctx context.Context, refreshToken string) error {
	query := `UPDATE sessions SET revoked_at = now() WHERE refresh_token = $1`
	_, err := r.db.Exec(ctx, query, refreshToken)
	return err
}

func (r *repository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *repository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Session, int64, error) {
	countQuery := `SELECT COUNT(*) FROM sessions WHERE user_id = $1 AND revoked_at IS NULL`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at, revoked_at, created_at
		FROM sessions
		WHERE user_id = $1 AND revoked_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.RefreshToken, &s.UserAgent, &s.IPAddress, &s.ExpiresAt, &s.RevokedAt, &s.CreatedAt); err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, s)
	}

	return sessions, total, nil
}

func (r *repository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM sessions WHERE expires_at < $1 OR revoked_at IS NOT NULL`
	result, err := r.db.Exec(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
