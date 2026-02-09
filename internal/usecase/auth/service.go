package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/auth"
	authShared "github.com/ap1-final-mini-moodle/internal/shared/auth"
	"github.com/google/uuid"
)

// Service для работы с JWT и refresh токенами
type Service struct {
	repo      auth.Repository
	jwtConfig *authShared.JWTConfig
}

func NewService(repo auth.Repository, jwtConfig *authShared.JWTConfig) *Service {
	return &Service{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

// IssueTokens создаёт новую сессию и выдаёт JWT + refresh token
func (s *Service) IssueTokens(ctx context.Context, user *authShared.UserData, userAgent, ipAddress string) (*auth.Token, error) {
	// Генерируем токены
	token, err := authShared.GenerateTokens(user, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Сохраняем сессию в БД
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

// RefreshAccessToken обновляет access token используя refresh token
func (s *Service) RefreshAccessToken(ctx context.Context, refreshToken string) (*auth.Token, error) {
	// Получаем сессию из БД
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Проверяем, не истекла ли сессия
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Здесь должны быть данные пользователя, получаем из БД по userID
	// В реальной системе здесь нужно получить данные пользователя
	// Для этого потребуется инъекция user repository или сохранение данных в refresh_sessions

	// Обновляем время последней активности
	if err := s.repo.UpdateSessionActivity(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Генерируем новый access token (refresh token остаётся тем же)
	userData := &authShared.UserData{
		ID:        session.UserID,
		Email:     "", // Нужно получить из БД
		FirstName: "",
		LastName:  "",
		Role:      "", // Нужно получить из БД
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

// VerifyAccessToken проверяет access token
func (s *Service) VerifyAccessToken(accessToken string) (*auth.TokenClaims, error) {
	claims, err := authShared.VerifyAccessToken(accessToken, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}
	return claims, nil
}

// RevokeSession деактивирует сессию (logout)
func (s *Service) RevokeSession(ctx context.Context, refreshToken string) error {
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	return s.repo.RevokeSession(ctx, session.ID)
}

// RevokeAllSessions деактивирует все сессии пользователя (logout everywhere)
func (s *Service) RevokeAllSessions(ctx context.Context, userID string) error {
	return s.repo.RevokeAllUserSessions(ctx, userID)
}

// GetActiveSessionsForUser получает все активные сессии пользователя
func (s *Service) GetActiveSessionsForUser(ctx context.Context, userID string) ([]*auth.RefreshSession, error) {
	return s.repo.GetActiveSessionsByUserID(ctx, userID)
}

// CleanupExpiredSessions удаляет истёкшие сессии
func (s *Service) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	return s.repo.DeleteExpiredSessions(ctx)
}
