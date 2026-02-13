package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AppealModule struct {
	handler     *AppealHandler
	authService *authUC.Service
}

func NewAppealModule(handler *AppealHandler, authService *authUC.Service) *AppealModule {
	return &AppealModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *AppealModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthTokenMiddleware(m.authService))

	appeals := api.Group("/appeals")
	{
		// VIEW — student видит свои, admin видит все
		appeals.GET("", m.handler.GetStudentAppeals)
		appeals.GET("/:id", m.handler.GetByID)

		// CREATE — только student
		appeals.POST("",
			middleware.RequireRole("student"),
			m.handler.Create,
		)

		// DELETE — только student (owner) или admin
		appeals.DELETE("/:id",
			middleware.RequireRole("student", "admin", "super_admin"),
			m.handler.Delete,
		)

		// RESOLVE — только teacher/admin
		appeals.PUT("/:id/resolve",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Resolve,
		)
	}

	api.GET("/teacher/appeals",
		middleware.RequireRole("teacher", "admin", "super_admin"),
		m.handler.GetTeacherAppeals,
	)
}
