package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type UploadModule struct {
	handler     *UploadHandler
	authService *authUC.Service
}

func NewUploadModule(handler *UploadHandler, authService *authUC.Service) *UploadModule {
	return &UploadModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *UploadModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	uploads := api.Group("/uploads")
	uploads.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		uploads.POST("", m.handler.Upload)
		uploads.GET("/:id", m.handler.GetByID)
		uploads.GET("/user/:userID", m.handler.ListByUser)
		uploads.DELETE("/:id", m.handler.Delete)
	}
}
