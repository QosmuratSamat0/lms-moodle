package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the session service interface
type Service interface {
	Create(ctx context.Context, userID uuid.UUID, userAgent, ipAddress *string, ttl time.Duration) (*Session, error)
	Validate(ctx context.Context, refreshToken string) (*SessionWithUser, error)
	Refresh(ctx context.Context, refreshToken string, userAgent, ipAddress *string, ttl time.Duration) (*Session, error)
	Revoke(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error
	RevokeByToken(ctx context.Context, refreshToken string) error
	RevokeAll(ctx context.Context, userID uuid.UUID) error
	ListMySessions(ctx context.Context, userID uuid.UUID, page, limit int) (*SessionListResponse, error)
	CleanupExpired(ctx context.Context) (int64, error)
}

type service struct {
	repo Repository
}

// NewService creates a new session service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, userAgent, ipAddress *string, ttl time.Duration) (*Session, error) {
	token, err := generateRefreshToken()
	if err != nil {
		return nil, errorx.Wrap(err, "generate refresh token")
	}

	session := &Session{
		UserID:       userID,
		RefreshToken: token,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(ttl),
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, errorx.Wrap(err, "create session")
	}

	return session, nil
}

func (s *service) Validate(ctx context.Context, refreshToken string) (*SessionWithUser, error) {
	session, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewUnauthorizedError("invalid refresh token")
		}
		return nil, errorx.Wrap(err, "get session")
	}

	if !session.IsValid() {
		return nil, errorx.NewUnauthorizedError("session expired or revoked")
	}

	if !session.UserIsActive {
		return nil, errorx.NewUnauthorizedError("user account is disabled")
	}

	return session, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string, userAgent, ipAddress *string, ttl time.Duration) (*Session, error) {
	// Validate existing session
	existing, err := s.Validate(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Revoke old session
	if err := s.repo.RevokeByRefreshToken(ctx, refreshToken); err != nil {
		return nil, errorx.Wrap(err, "revoke old session")
	}

	// Create new session
	return s.Create(ctx, existing.UserID, userAgent, ipAddress, ttl)
}

func (s *service) Revoke(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error {
	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("session")
		}
		return errorx.Wrap(err, "get session")
	}

	// Verify ownership
	if session.UserID != userID {
		return errorx.NewForbiddenError("you can only revoke your own sessions")
	}

	if err := s.repo.RevokeByID(ctx, sessionID); err != nil {
		return errorx.Wrap(err, "revoke session")
	}

	return nil
}

func (s *service) RevokeByToken(ctx context.Context, refreshToken string) error {
	if err := s.repo.RevokeByRefreshToken(ctx, refreshToken); err != nil {
		return errorx.Wrap(err, "revoke session")
	}
	return nil
}

func (s *service) RevokeAll(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.RevokeAllByUserID(ctx, userID); err != nil {
		return errorx.Wrap(err, "revoke all sessions")
	}
	return nil
}

func (s *service) ListMySessions(ctx context.Context, userID uuid.UUID, page, limit int) (*SessionListResponse, error) {
	offset := (page - 1) * limit
	sessions, total, err := s.repo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list sessions")
	}

	var responses []SessionResponse
	for _, sess := range sessions {
		responses = append(responses, SessionResponse{
			ID:        sess.ID.String(),
			UserAgent: sess.UserAgent,
			IPAddress: sess.IPAddress,
			ExpiresAt: sess.ExpiresAt.Format(time.RFC3339),
			CreatedAt: sess.CreatedAt.Format(time.RFC3339),
		})
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &SessionListResponse{
		Sessions:   responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *service) CleanupExpired(ctx context.Context) (int64, error) {
	count, err := s.repo.DeleteExpired(ctx)
	if err != nil {
		return 0, errorx.Wrap(err, "cleanup expired sessions")
	}
	return count, nil
}

// Helper function to generate refresh token
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

var _ Service = (*service)(nil)

// Response types

// SessionResponse represents a session in API responses
type SessionResponse struct {
	ID        string  `json:"id"`
	UserAgent *string `json:"user_agent,omitempty"`
	IPAddress *string `json:"ip_address,omitempty"`
	ExpiresAt string  `json:"expires_at"`
	CreatedAt string  `json:"created_at"`
}

// SessionListResponse represents paginated list of sessions
type SessionListResponse struct {
	Sessions   []SessionResponse `json:"sessions"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}
