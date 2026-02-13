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

	// Public route for serving uploaded files (no auth needed)
	api.GET("/uploads/files/:filename", m.handler.ServeFile)

	uploads := api.Group("/uploads")
	uploads.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// VIEW — owner или admin
		uploads.GET("/:id", m.handler.GetByID)
		uploads.GET("/user/:userID", m.handler.ListByUser)

		// CREATE — все могут загружать
		uploads.POST("",
			m.handler.Upload,
		)

		// DELETE — owner или admin
		uploads.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
