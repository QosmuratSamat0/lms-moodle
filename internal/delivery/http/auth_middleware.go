package http

import (
	"fmt"

	"github.com/ap1-final-mini-moodle/internal/domain/auth"
	"github.com/gin-gonic/gin"
)

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
