// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/ap1-final-mini-moodle/internal/domain/enrollment"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// EnrollmentRouter handles enrollment-related routes
type EnrollmentRouter struct {
	enrollmentHandler *enrollment.Handler
}

// NewEnrollmentRouter creates a new enrollment router
func NewEnrollmentRouter(enrollmentHandler *enrollment.Handler) *EnrollmentRouter {
	return &EnrollmentRouter{
		enrollmentHandler: enrollmentHandler,
	}
}

// SetupRoutes configures enrollment routes
func (er *EnrollmentRouter) SetupRoutes(api *gin.RouterGroup) {
	enrollments := api.Group("/enrollments")
	{
		// Student enrollment
		enrollments.POST("", er.enrollmentHandler.Enroll)
		enrollments.DELETE("/:id", er.enrollmentHandler.Drop)
		enrollments.GET("/my", er.enrollmentHandler.ListMyEnrollments)

		// Course enrollments (teachers)
		api.GET("/courses/:id/enrollments", er.enrollmentHandler.ListByCourse)

		// Approval (teacher/admin)
		teacherRoutes := enrollments.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin", "manager"))
		{
			teacherRoutes.PUT("/:id/status", er.enrollmentHandler.UpdateStatus)
		}
	}
}
