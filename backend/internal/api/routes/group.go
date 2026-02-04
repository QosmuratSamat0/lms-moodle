// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/ap1-final-mini-moodle/internal/domain/group"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// GroupRouter handles group-related routes
type GroupRouter struct {
	groupHandler *group.Handler
}

// NewGroupRouter creates a new group router
func NewGroupRouter(groupHandler *group.Handler) *GroupRouter {
	return &GroupRouter{
		groupHandler: groupHandler,
	}
}

// SetupRoutes configures group routes
func (gr *GroupRouter) SetupRoutes(api *gin.RouterGroup) {
	groups := api.Group("/groups")
	{
		// View groups (all authenticated users)
		groups.GET("", gr.groupHandler.List)
		groups.GET("/:id", gr.groupHandler.GetByID)
		groups.GET("/code/:code", gr.groupHandler.GetByCode)
		groups.GET("/:id/assignments", gr.groupHandler.ListGroupAssignments)

		// Admin/Manager only
		adminRoutes := groups.Group("")
		adminRoutes.Use(middleware.RequireRole("admin", "manager"))
		{
			adminRoutes.POST("", gr.groupHandler.Create)
			adminRoutes.PUT("/:id", gr.groupHandler.Update)
			adminRoutes.DELETE("/:id", gr.groupHandler.Delete)
			adminRoutes.POST("/assignments", gr.groupHandler.AssignTeacher)
			adminRoutes.DELETE("/assignments/:id", gr.groupHandler.UnassignTeacher)
		}
	}
}
