// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/attendance"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// AttendanceRouter handles attendance-related routes
type AttendanceRouter struct {
	attendanceHandler *attendance.Handler
}

// NewAttendanceRouter creates a new attendance router
func NewAttendanceRouter(attendanceHandler *attendance.Handler) *AttendanceRouter {
	return &AttendanceRouter{
		attendanceHandler: attendanceHandler,
	}
}

// SetupRoutes configures attendance routes
func (ar *AttendanceRouter) SetupRoutes(api *gin.RouterGroup) {
	attendance := api.Group("/attendance")
	{
		// Sessions
		attendance.GET("/sessions/:id", ar.attendanceHandler.GetSession)
		api.GET("/courses/:id/attendance/sessions", ar.attendanceHandler.ListSessionsByCourse)

		// Teacher/Admin only
		teacherRoutes := attendance.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("/sessions", ar.attendanceHandler.CreateSession)
			teacherRoutes.PUT("/sessions/:id", ar.attendanceHandler.UpdateSession)
			teacherRoutes.DELETE("/sessions/:id", ar.attendanceHandler.DeleteSession)
			teacherRoutes.POST("/sessions/:id/mark", ar.attendanceHandler.MarkAttendance)
			teacherRoutes.POST("/sessions/:id/bulk-mark", ar.attendanceHandler.BulkMarkAttendance)
		}

		// Student attendance summary (moved under /attendance to avoid conflict)
		attendance.GET("/students/:student_id/courses/:course_id/summary", ar.attendanceHandler.GetStudentCourseSummary)
	}
}
