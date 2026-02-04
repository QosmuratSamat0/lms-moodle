// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/ap1-final-mini-moodle/internal/domain/manager"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// ManagerRouter handles manager-related routes
type ManagerRouter struct {
	managerHandler *manager.Handler
}

// NewManagerRouter creates a new manager router
func NewManagerRouter(managerHandler *manager.Handler) *ManagerRouter {
	return &ManagerRouter{
		managerHandler: managerHandler,
	}
}

// SetupRoutes configures manager routes
func (mr *ManagerRouter) SetupRoutes(api *gin.RouterGroup) {
	managers := api.Group("/managers")
	managers.Use(middleware.RequireRole("manager", "admin"))
	{
		managers.GET("/me", mr.managerHandler.GetProfile)
		managers.PUT("/me", mr.managerHandler.UpdateProfile)
		managers.GET("/overview", mr.managerHandler.GetSystemOverview)

		// User management
		managers.POST("/users/:id/activate", mr.managerHandler.ActivateUser)
		managers.POST("/users/:id/deactivate", mr.managerHandler.DeactivateUser)
		managers.POST("/users/bulk", mr.managerHandler.BulkUserAction)

		// Admin only
		adminRoutes := managers.Group("")
		adminRoutes.Use(middleware.RequireRole("admin"))
		{
			adminRoutes.GET("", mr.managerHandler.List)
			adminRoutes.GET("/:id", mr.managerHandler.GetByID)
		}
	}
}
