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
		appeals.POST("", m.handler.Create)
		appeals.GET("", m.handler.GetStudentAppeals)
		appeals.GET("/:id", m.handler.GetByID)
		appeals.DELETE("/:id", m.handler.Delete)
		appeals.PUT("/:id/resolve", m.handler.Resolve)
	}

	api.GET("/teacher/appeals", m.handler.GetTeacherAppeals)
}
