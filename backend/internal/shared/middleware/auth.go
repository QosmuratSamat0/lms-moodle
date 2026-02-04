// internal/shared/middleware/auth.go
package middleware

import (
	"strings"

	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
	RoleKey   contextKey = "role"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtManager *utils.JWTManager
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtManager *utils.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// Authenticate is a Gin middleware that validates JWT tokens
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// First check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token = parts[1]
			}
		}

		// Fallback to query parameter for WebSocket connections
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			utils.Unauthorized(c, "missing authorization")
			c.Abort()
			return
		}

		claims, err := m.jwtManager.ValidateAccessToken(token)
		if err != nil {
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(string(UserIDKey), claims.UserID)
		c.Set(string(EmailKey), claims.Email)
		c.Set(string(RoleKey), claims.Role)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}
		utils.Forbidden(c, "insufficient permissions")
		c.Abort()
	}
}

// RequireAnyRole middleware checks if user has any of the required roles
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return RequireRole(roles...)
}

// GetUserID extracts user ID from Gin context
// Returns uuid.Nil if not found
func GetUserID(c *gin.Context) uuid.UUID {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		return uuid.Nil
	}
	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}

// MustGetUserID extracts user ID from Gin context with error handling
func MustGetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetUserEmail extracts user email from Gin context
func GetUserEmail(c *gin.Context) string {
	email, exists := c.Get(string(EmailKey))
	if !exists {
		return ""
	}
	e, ok := email.(string)
	if !ok {
		return ""
	}
	return e
}

// GetUserRole extracts user role from Gin context
func GetUserRole(c *gin.Context) string {
	role, exists := c.Get(string(RoleKey))
	if !exists {
		return ""
	}
	r, ok := role.(string)
	if !ok {
		return ""
	}
	return r
}
