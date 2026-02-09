package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AuthModule struct {
	handler     *AuthHandler
	authService *authUC.Service
	secret      []byte
}

func NewAuthModule(handler *AuthHandler, authService *authUC.Service, secret []byte) *AuthModule {
	return &AuthModule{
		handler:     handler,
		authService: authService,
		secret:      secret,
	}
}

func (m *AuthModule) Register(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	{
		// Публичные endpoints
		auth.POST("/login", m.handler.Login)
		auth.POST("/refresh", m.handler.RefreshToken)
		auth.POST("/verify", m.handler.Verify)

		// Private endpoints (требуют аутентификации)
		protected := auth.Group("")
		protected.Use(middleware.AuthTokenMiddleware(m.authService))
		{
			protected.POST("/logout", m.handler.Logout)
			protected.POST("/logout-everywhere", m.handler.LogoutEverywhere)
			protected.GET("/sessions", m.handler.GetActiveSessions)
			protected.DELETE("/sessions/:session_id", m.handler.RevokeSessionByID)
		}
	}
}
