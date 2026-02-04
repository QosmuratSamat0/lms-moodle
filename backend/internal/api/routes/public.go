// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/user"
	"github.com/gin-gonic/gin"
)

// PublicRouter handles public routes (no authentication required)
type PublicRouter struct {
	userHandler *user.Handler
}

// NewPublicRouter creates a new public router
func NewPublicRouter(userHandler *user.Handler) *PublicRouter {
	return &PublicRouter{
		userHandler: userHandler,
	}
}

// SetupRoutes configures public routes
func (pr *PublicRouter) SetupRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", pr.userHandler.Register)
		auth.POST("/login", pr.userHandler.Login)
		auth.POST("/refresh", pr.userHandler.RefreshToken)
	}
}
