package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type CourseCategoryModule struct {
	handler     *CourseCategoryHandler
	authService *authUC.Service
}

func NewCourseCategoryModule(handler *CourseCategoryHandler, authService *authUC.Service) *CourseCategoryModule {
	return &CourseCategoryModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *CourseCategoryModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	categories := api.Group("/categories")
	categories.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// VIEW — все могут видеть категории
		categories.GET("", m.handler.List)
		categories.GET("/:id", m.handler.GetByID)

		// CREATE — только admin
		categories.POST("",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Create,
		)

		// UPDATE — only admin
		categories.PUT("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Update,
		)

		// DELETE — только admin
		categories.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
