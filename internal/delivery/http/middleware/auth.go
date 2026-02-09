package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ap1-final-mini-moodle/internal/domain/auth"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

// AuthTokenMiddleware извлекает и проверяет JWT токен из заголовка Authorization
func AuthTokenMiddleware(authService *authUC.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Ожидаем формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		accessToken := parts[1]

		// Проверяем токен
		claims, err := authService.VerifyAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Сохраняем данные в контекст
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		c.Set("refreshToken", claims.RefreshToken)
		c.Set("tokenClaims", claims)

		c.Next()
	}
}

// OptionalAuthTokenMiddleware опциональная аутентификация (не требует токена, но использует его если есть)
func OptionalAuthTokenMiddleware(authService *authUC.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Токена нет, но это нормально
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Неправильный формат, но не критично
			c.Next()
			return
		}

		accessToken := parts[1]

		// Проверяем токен
		claims, err := authService.VerifyAccessToken(accessToken)
		if err != nil {
			// Токен невалиден, но это не критично для опциональной аутентификации
			c.Next()
			return
		}

		// Сохраняем данные в контекст
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		c.Set("refreshToken", claims.RefreshToken)
		c.Set("tokenClaims", claims)

		c.Next()
	}
}

// GetUserIDFromContext безопасно получает userID из контекста
func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", fmt.Errorf("user not authenticated")
	}

	id, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("invalid user id type")
	}

	if id == "" {
		return "", fmt.Errorf("user id is empty")
	}

	return id, nil
}

// GetTokenClaimsFromContext получает полные claims из контекста
func GetTokenClaimsFromContext(c *gin.Context) (*auth.TokenClaims, error) {
	claims, exists := c.Get("tokenClaims")
	if !exists {
		return nil, fmt.Errorf("token claims not found")
	}

	tokenClaims, ok := claims.(*auth.TokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims type")
	}

	return tokenClaims, nil
}
