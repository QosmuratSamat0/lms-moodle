// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/ap1-final-mini-moodle/internal/domain/assignment"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// AssignmentRouter handles assignment-related routes
type AssignmentRouter struct {
	assignmentHandler *assignment.Handler
}

// NewAssignmentRouter creates a new assignment router
func NewAssignmentRouter(assignmentHandler *assignment.Handler) *AssignmentRouter {
	return &AssignmentRouter{
		assignmentHandler: assignmentHandler,
	}
}

// SetupRoutes configures assignment routes
func (ar *AssignmentRouter) SetupRoutes(api *gin.RouterGroup) {
	assignments := api.Group("/assignments")
	{
		assignments.GET("/:id", ar.assignmentHandler.GetByID)

		// Course assignments
		api.GET("/courses/:id/assignments", ar.assignmentHandler.ListByCourse)

		// Teacher/Admin only
		teacherRoutes := assignments.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("", ar.assignmentHandler.Create)
			teacherRoutes.PUT("/:id", ar.assignmentHandler.Update)
			teacherRoutes.DELETE("/:id", ar.assignmentHandler.Delete)
		}
	}
}
