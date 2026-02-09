package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type GradeModule struct {
	handler     *GradeHandler
	authService *authUC.Service
}

func NewGradeModule(handler *GradeHandler, authService *authUC.Service) *GradeModule {
	return &GradeModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *GradeModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	grades := api.Group("/grades")
	grades.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// VIEW — student видит свои, teacher видит свои, admin видит все
		grades.GET("/:id", m.handler.GetByID)

		// CREATE — только teacher/admin
		grades.POST("",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Grade,
		)

		// DELETE — только admin
		grades.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
