// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/ap1-final-mini-moodle/internal/domain/course"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// CourseRouter handles course-related routes
type CourseRouter struct {
	courseHandler *course.Handler
}

// NewCourseRouter creates a new course router
func NewCourseRouter(courseHandler *course.Handler) *CourseRouter {
	return &CourseRouter{
		courseHandler: courseHandler,
	}
}

// SetupRoutes configures course routes
func (cr *CourseRouter) SetupRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses")
	{
		courses.GET("", cr.courseHandler.List)
		courses.GET("/:id", cr.courseHandler.GetByID)

		// Teacher/Admin only
		teacherRoutes := courses.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("", cr.courseHandler.Create)
			teacherRoutes.PUT("/:id", cr.courseHandler.Update)
			teacherRoutes.DELETE("/:id", cr.courseHandler.Delete)
		}
	}
}
