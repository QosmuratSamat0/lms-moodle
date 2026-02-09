package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type DashboardModule struct {
	handler     *DashboardHandler
	authService *authUC.Service
}

func NewDashboardModule(handler *DashboardHandler, authService *authUC.Service) *DashboardModule {
	return &DashboardModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *DashboardModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// Student dashboard — только для студентов
		dashboard.GET("/student",
			middleware.RequireRole("student"),
			m.handler.GetStudentDashboard,
		)
		dashboard.GET("/student/courses/:courseId",
			middleware.RequireRole("student"),
			m.handler.GetStudentCourseStats,
		)

		// Teacher dashboard — только для учителей/админов
		dashboard.GET("/teacher",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.GetTeacherDashboard,
		)
	}
}
