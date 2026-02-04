// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/student"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// StudentRouter handles student-related routes
type StudentRouter struct {
	studentHandler *student.Handler
}

// NewStudentRouter creates a new student router
func NewStudentRouter(studentHandler *student.Handler) *StudentRouter {
	return &StudentRouter{
		studentHandler: studentHandler,
	}
}

// SetupRoutes configures student routes
func (sr *StudentRouter) SetupRoutes(api *gin.RouterGroup) {
	students := api.Group("/students")
	{
		// Student self-management
		studentRoutes := students.Group("")
		studentRoutes.Use(middleware.RequireRole("student"))
		{
			studentRoutes.GET("/me", sr.studentHandler.GetProfile)
			studentRoutes.PUT("/me", sr.studentHandler.UpdateProfile)
			studentRoutes.GET("/me/stats", sr.studentHandler.GetStats)
		}

		// View students (teachers, managers)
		viewRoutes := students.Group("")
		viewRoutes.Use(middleware.RequireRole("teacher", "admin", "manager"))
		{
			viewRoutes.GET("", sr.studentHandler.List)
			viewRoutes.GET("/:id", sr.studentHandler.GetByID)
			viewRoutes.GET("/:id/stats", sr.studentHandler.GetStatsByID)
			viewRoutes.GET("/group/:group", sr.studentHandler.ListByGroup)
			viewRoutes.GET("/course/:id", sr.studentHandler.ListByCourse)
		}
	}
}
