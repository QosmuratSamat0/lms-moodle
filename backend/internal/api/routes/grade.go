// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/grade"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// GradeRouter handles grade-related routes
type GradeRouter struct {
	gradeHandler *grade.Handler
}

// NewGradeRouter creates a new grade router
func NewGradeRouter(gradeHandler *grade.Handler) *GradeRouter {
	return &GradeRouter{
		gradeHandler: gradeHandler,
	}
}

// SetupRoutes configures grade routes
func (gr *GradeRouter) SetupRoutes(api *gin.RouterGroup) {
	grades := api.Group("/grades")
	{
		grades.GET("/:id", gr.gradeHandler.GetByID)
		grades.GET("/my", gr.gradeHandler.ListMyGrades)

		// Teacher grading
		teacherRoutes := grades.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("", gr.gradeHandler.GradeSubmission)
			teacherRoutes.PUT("/:id", gr.gradeHandler.Update)
			teacherRoutes.DELETE("/:id", gr.gradeHandler.Delete)
		}

		// Course grades
		api.GET("/courses/:id/grades", gr.gradeHandler.ListByCourse)
	}
}
