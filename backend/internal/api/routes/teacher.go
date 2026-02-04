// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/group"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/teacher"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// TeacherRouter handles teacher-related routes
type TeacherRouter struct {
	teacherHandler *teacher.Handler
	groupHandler   *group.Handler
}

// NewTeacherRouter creates a new teacher router
func NewTeacherRouter(teacherHandler *teacher.Handler, groupHandler *group.Handler) *TeacherRouter {
	return &TeacherRouter{
		teacherHandler: teacherHandler,
		groupHandler:   groupHandler,
	}
}

// SetupRoutes configures teacher routes
func (tr *TeacherRouter) SetupRoutes(api *gin.RouterGroup) {
	teachers := api.Group("/teachers")
	{
		// Teacher self-management (allow admin too for testing)
		teacherRoutes := teachers.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.GET("/me", tr.teacherHandler.GetProfile)
			teacherRoutes.PUT("/me", tr.teacherHandler.UpdateProfile)
			teacherRoutes.GET("/me/stats", tr.teacherHandler.GetStats)
			teacherRoutes.GET("/me/group-assignments", tr.groupHandler.ListTeacherAssignments) // Teacher's group assignments
		}

		// View teachers (all authenticated)
		teachers.GET("", tr.teacherHandler.List)
		teachers.GET("/:id", tr.teacherHandler.GetByID)

		// Admin/Manager only
		adminRoutes := teachers.Group("")
		adminRoutes.Use(middleware.RequireRole("admin", "manager"))
		{
			adminRoutes.GET("/:id/stats", tr.teacherHandler.GetStatsByID)
			adminRoutes.GET("/:id/group-assignments", tr.groupHandler.ListTeacherAssignmentsByID) // View specific teacher's assignments
			adminRoutes.GET("/department/:department", tr.teacherHandler.ListByDepartment)
		}
	}
}
