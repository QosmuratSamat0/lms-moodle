package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ap1-final-mini-moodle/internal/domain/auth"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

func AuthTokenMiddleware(authService *authUC.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("[AUTH] Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		log.Printf("[AUTH] Authorization header received: %s (length: %d)", authHeader[:min(len(authHeader), 50)], len(authHeader))

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			log.Printf("[AUTH] Invalid format: got %d parts instead of 2", len(parts))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		if parts[0] != "Bearer" {
			log.Printf("[AUTH] Invalid scheme: got '%s' instead of 'Bearer'", parts[0])
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		accessToken := parts[1]
		log.Printf("[AUTH] Token extracted (length: %d)", len(accessToken))

		claims, err := authService.VerifyAccessToken(accessToken)
		if err != nil {
			log.Printf("[AUTH] Token verification failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token", "details": err.Error()})
			c.Abort()
			return
		}

		log.Printf("[AUTH] Token verified for user: %s (role: %s)", claims.UserID, claims.Role)
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

func OptionalAuthTokenMiddleware(authService *authUC.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		accessToken := parts[1]

		claims, err := authService.VerifyAccessToken(accessToken)
		if err != nil {
			c.Next()
			return
		}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
