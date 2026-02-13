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
	att := api.Group("/attendance")
	att.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// Sessions
		att.GET("/course/:courseId/sessions", m.handler.ListSessionsByCourse)
		att.GET("/sessions/:id", m.handler.GetSession)
		att.POST("/sessions",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.CreateSession,
		)
		att.DELETE("/sessions/:id",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.DeleteSession,
		)

		// Marks
		att.GET("/marks/session/:sessionId", m.handler.GetMarksBySession)
		att.POST("/marks/bulk",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.BulkMark,
		)

		// Student view
		att.GET("/course/:courseId/me", m.handler.GetMyAttendance)
		att.GET("/course/:courseId/student/:studentId", m.handler.GetStudentSummary)
	}
}
