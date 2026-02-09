package auth

import (
	"context"
	"time"
)

// Token представляет JWT токен с метаданными
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"` // в секундах
	IssuedAt     time.Time `json:"-"`
}

// TokenClaims содержит данные, кодируемые в JWT
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

// RefreshSession хранит информацию о сессии с refresh токеном
type RefreshSession struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	IssuedAt     time.Time `db:"issued_at"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	// Для отслеживания устройств
	UserAgent string `db:"user_agent"`
	IPAddress string `db:"ip_address"`
}

// CreateSessionInput для создания новой успешной сессии
type CreateSessionInput struct {
	UserID       string
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

// Repository интерфейс для работы с refresh токенами
type Repository interface {
	// CreateSession создаёт новую сессию с refresh токеном
	CreateSession(ctx context.Context, session *RefreshSession) error

	// GetSessionByRefreshToken получает сессию по refresh токену
	GetSessionByRefreshToken(ctx context.Context, token string) (*RefreshSession, error)

	// GetSessionByID получает сессию по ID
	GetSessionByID(ctx context.Context, id string) (*RefreshSession, error)

	// GetActiveSessionsByUserID получает все активные сессии пользователя
	GetActiveSessionsByUserID(ctx context.Context, userID string) ([]*RefreshSession, error)

	// UpdateSessionActivity обновляет время последней активности
	UpdateSessionActivity(ctx context.Context, id string) error

	// RevokeSession деактивирует сессию (logout)
	RevokeSession(ctx context.Context, id string) error

	// RevokeAllUserSessions деактивирует все сессии пользователя (logout from everywhere)
	RevokeAllUserSessions(ctx context.Context, userID string) error

	// DeleteExpiredSessions удаляет истёкшие сессии (для очистки БД)
	DeleteExpiredSessions(ctx context.Context) (int64, error)
}
