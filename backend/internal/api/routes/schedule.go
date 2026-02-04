// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/ap1-final-mini-moodle/internal/domain/schedule"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// ScheduleRouter handles schedule-related routes
type ScheduleRouter struct {
	scheduleHandler *schedule.Handler
}

// NewScheduleRouter creates a new schedule router
func NewScheduleRouter(scheduleHandler *schedule.Handler) *ScheduleRouter {
	return &ScheduleRouter{
		scheduleHandler: scheduleHandler,
	}
}

// SetupRoutes configures schedule routes
func (sr *ScheduleRouter) SetupRoutes(api *gin.RouterGroup) {
	schedule := api.Group("/schedule")
	{
		schedule.GET("", sr.scheduleHandler.ListMySchedule)
		schedule.GET("/calendar", sr.scheduleHandler.GetCalendar)
		schedule.GET("/events/:id", sr.scheduleHandler.GetByID)

		// Course schedule
		api.GET("/courses/:id/schedule", sr.scheduleHandler.ListByCourse)

		// Teacher/Admin only
		teacherRoutes := schedule.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.POST("/events", sr.scheduleHandler.Create)
			teacherRoutes.PUT("/events/:id", sr.scheduleHandler.Update)
			teacherRoutes.DELETE("/events/:id", sr.scheduleHandler.Delete)
		}
	}
}
