package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authShared "github.com/ap1-final-mini-moodle/internal/shared/auth"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	userUC "github.com/ap1-final-mini-moodle/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *userUC.Service
	authService *authUC.Service
}

func NewAuthHandler(userService *userUC.Service, authService *authUC.Service) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
	}
}

// Login аутентифицирует пользователя и выдаёт JWT + refresh token
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем пользователя
	user, err := h.userService.GetByEmailAndPassword(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Подготавливаем данные пользователя для токена
	userData := &authShared.UserData{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
	}

	// Получаем user agent и IP адрес
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	// Выдаём токены
	token, err := h.authService.IssueTokens(c.Request.Context(), userData, userAgent, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expires_in":    token.ExpiresIn,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
		},
	})
}

// RefreshToken обновляет access token, используя refresh token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем новые токены
	token, err := h.authService.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expires_in":    token.ExpiresIn,
	})
}

// Logout деактивирует текущую сессию
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
		return
	}

	if err := h.authService.RevokeSession(c.Request.Context(), req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// LogoutEverywhere деактивирует все сессии пользователя
func (h *AuthHandler) LogoutEverywhere(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.authService.RevokeAllSessions(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout from all devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out from all devices successfully"})
}

// GetActiveSessions получает список активных сессий пользователя
func (h *AuthHandler) GetActiveSessions(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sessions, err := h.authService.GetActiveSessionsForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sessions"})
		return
	}

	// Преобразуем для API
	var sessionsData []gin.H
	for _, session := range sessions {
		sessionsData = append(sessionsData, gin.H{
			"id":         session.ID,
			"user_agent": session.UserAgent,
			"ip_address": session.IPAddress,
			"issued_at":  session.IssuedAt,
			"expires_at": session.ExpiresAt,
			"created_at": session.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessionsData,
		"count":    len(sessions),
	})
}

// RevokeSessionByID деактивирует конкретную сессию по ID
func (h *AuthHandler) RevokeSessionByID(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Проверяем, что это сессия пользователя (безопасность)
	sessions, err := h.authService.GetActiveSessionsForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify session"})
		return
	}

	found := false
	for _, session := range sessions {
		if session.ID == sessionID {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusForbidden, gin.H{"error": "this session does not belong to you"})
		return
	}

	// Используем repository напрямую (нужно добавить в auth service)
	// Временно возвращаем ошибку
	c.JSON(http.StatusInternalServerError, gin.H{"error": "not implemented"})
}

// Verify проверяет валидность текущего токена
func (h *AuthHandler) Verify(c *gin.Context) {
	claims, err := GetTokenClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    claims.UserID,
		"email":      claims.Email,
		"role":       claims.Role,
		"first_name": claims.FirstName,
		"last_name":  claims.LastName,
		"issued_at":  claims.IssuedAt,
		"expires_at": claims.ExpiresAt,
	})
}
