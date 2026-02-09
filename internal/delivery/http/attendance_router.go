package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AttendanceModule struct {
	handler     *AttendanceHandler
	authService *authUC.Service
}

func NewAttendanceModule(handler *AttendanceHandler, authService *authUC.Service) *AttendanceModule {
	return &AttendanceModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *AttendanceModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	attendance := api.Group("/attendance")
	attendance.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// VIEW — teacher видит свои, admin видит все
		attendance.GET("/course/:courseId", m.handler.ListByCourse)

		// CREATE — только teacher/admin
		attendance.POST("",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Record,
		)

		// DELETE — только admin
		attendance.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
