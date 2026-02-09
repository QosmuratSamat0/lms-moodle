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
		categories.POST("", m.handler.Create)
		categories.GET("", m.handler.List)
		categories.GET("/:id", m.handler.GetByID)
		categories.PUT("/:id", m.handler.Update)
		categories.DELETE("/:id", m.handler.Delete)
	}
}
