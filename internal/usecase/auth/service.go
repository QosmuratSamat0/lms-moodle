package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/auth"
	"github.com/ap1-final-mini-moodle/internal/domain/user"
	authShared "github.com/ap1-final-mini-moodle/internal/shared/auth"
	"github.com/google/uuid"
)


type Service struct {
	repo      auth.Repository
	userRepo  user.Repository
	jwtConfig *authShared.JWTConfig
}

func NewService(repo auth.Repository, userRepo user.Repository, jwtConfig *authShared.JWTConfig) *Service {
	return &Service{
		repo:      repo,
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

func (s *Service) IssueTokens(ctx context.Context, user *authShared.UserData, userAgent, ipAddress string) (*auth.Token, error) {
	token, err := authShared.GenerateTokens(user, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	session := &auth.RefreshSession{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    time.Now().Add(s.jwtConfig.RefreshTokenDuration),
		IssuedAt:     time.Now(),
		IsActive:     true,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return token, nil
}

func (s *Service) RefreshAccessToken(ctx context.Context, refreshToken string) (*auth.Token, error) {
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	usr, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := s.repo.UpdateSessionActivity(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	userData := &authShared.UserData{
		ID:        usr.ID,
		Email:     usr.Email,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Role:      string(usr.Role),
	}

	accessToken, err := authShared.GenerateAccessToken(userData, refreshToken, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &auth.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtConfig.AccessTokenDuration.Seconds()),
		IssuedAt:     time.Now(),
	}, nil
}

func (s *Service) VerifyAccessToken(accessToken string) (*auth.TokenClaims, error) {
	claims, err := authShared.VerifyAccessToken(accessToken, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}
	return claims, nil
}

func (s *Service) RevokeSession(ctx context.Context, refreshToken string) error {
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	return s.repo.RevokeSession(ctx, session.ID)
}

func (s *Service) RevokeAllSessions(ctx context.Context, userID string) error {
	return s.repo.RevokeAllUserSessions(ctx, userID)
}

func (s *Service) GetActiveSessionsForUser(ctx context.Context, userID string) ([]*auth.RefreshSession, error) {
	return s.repo.GetActiveSessionsByUserID(ctx, userID)
}

func (s *Service) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	return s.repo.DeleteExpiredSessions(ctx)
}
