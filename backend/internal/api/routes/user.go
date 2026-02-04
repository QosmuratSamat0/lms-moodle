// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/ap1-final-mini-moodle/internal/domain/user"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// UserRouter handles user-related routes
type UserRouter struct {
	userHandler *user.Handler
}

// NewUserRouter creates a new user router
func NewUserRouter(userHandler *user.Handler) *UserRouter {
	return &UserRouter{
		userHandler: userHandler,
	}
}

// SetupRoutes configures user routes
func (ur *UserRouter) SetupRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/me", ur.userHandler.GetProfile)
		users.PUT("/me", ur.userHandler.UpdateProfile)
		users.PUT("/me/password", ur.userHandler.ChangePassword)

		// Admin only
		admin := users.Group("")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("", ur.userHandler.ListUsers)
			admin.GET("/:id", ur.userHandler.GetUser)
			admin.POST("/:id/deactivate", ur.userHandler.DeactivateUser)
			admin.POST("/:id/activate", ur.userHandler.ActivateUser)
		}
	}
}
