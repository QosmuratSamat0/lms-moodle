// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/analytics"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// AnalyticsRouter handles analytics-related routes
type AnalyticsRouter struct {
	analyticsHandler *analytics.Handler
}

// NewAnalyticsRouter creates a new analytics router
func NewAnalyticsRouter(analyticsHandler *analytics.Handler) *AnalyticsRouter {
	return &AnalyticsRouter{
		analyticsHandler: analyticsHandler,
	}
}

// SetupRoutes configures analytics routes
func (ar *AnalyticsRouter) SetupRoutes(api *gin.RouterGroup) {
	analyticsRoutes := api.Group("/analytics")
	{
		// Track events (all users)
		analyticsRoutes.POST("/events", ar.analyticsHandler.TrackEvent)

		// Admin/Manager only
		adminRoutes := analyticsRoutes.Group("")
		adminRoutes.Use(middleware.RequireRole("admin", "manager"))
		{
			adminRoutes.GET("/system", ar.analyticsHandler.GetSystemStats)
		}

		// Course analytics (teachers)
		teacherRoutes := analyticsRoutes.Group("")
		teacherRoutes.Use(middleware.RequireRole("teacher", "admin"))
		{
			teacherRoutes.GET("/courses/:id", ar.analyticsHandler.GetCourseStats)
			teacherRoutes.GET("/courses/:id/leaderboard", ar.analyticsHandler.GetCourseLeaderboard)
		}
	}
}
