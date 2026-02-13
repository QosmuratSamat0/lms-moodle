package auth

import (
	"context"
	"time"
)

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"` 
	IssuedAt     time.Time `json:"-"`
}

type TokenClaims struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Role         string `json:"role"`
	IssuedAt     int64  `json:"iat"`
	ExpiresAt    int64  `json:"exp"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshSession struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	IssuedAt     time.Time `db:"issued_at"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	UserAgent string `db:"user_agent"`
	IPAddress string `db:"ip_address"`
}

type CreateSessionInput struct {
	UserID       string
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

type Repository interface {
	CreateSession(ctx context.Context, session *RefreshSession) error
	GetSessionByRefreshToken(ctx context.Context, token string) (*RefreshSession, error)
	GetSessionByID(ctx context.Context, id string) (*RefreshSession, error)
	GetActiveSessionsByUserID(ctx context.Context, userID string) ([]*RefreshSession, error)
	UpdateSessionActivity(ctx context.Context, id string) error
	RevokeSession(ctx context.Context, id string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	DeleteExpiredSessions(ctx context.Context) (int64, error)
}
