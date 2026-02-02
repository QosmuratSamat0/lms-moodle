package session

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session (for refresh tokens)
type Session struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	RefreshToken string     `json:"refresh_token" db:"refresh_token"`
	UserAgent    *string    `json:"user_agent" db:"user_agent"`
	IPAddress    *string    `json:"ip_address" db:"ip_address"`
	ExpiresAt    time.Time  `json:"expires_at" db:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at" db:"revoked_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// IsValid checks if session is valid (not expired, not revoked)
func (s *Session) IsValid() bool {
	if s.RevokedAt != nil {
		return false
	}
	return time.Now().Before(s.ExpiresAt)
}

// SessionWithUser includes user info
type SessionWithUser struct {
	Session
	UserEmail    string `json:"user_email" db:"user_email"`
	UserRole     string `json:"user_role" db:"user_role"`
	UserIsActive bool   `json:"user_is_active" db:"user_is_active"`
}
