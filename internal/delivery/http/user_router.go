package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type UserModule struct {
	handler     *UserHandler
	authService *authUC.Service
}

func NewUserModule(handler *UserHandler, authService *authUC.Service) *UserModule {
	return &UserModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *UserModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	users := api.Group("/users")
	{
		// Public
		users.POST("/register", m.handler.Register)
		users.POST("/login", m.handler.Login)

		// Protected
		protected := users.Group("")
		protected.Use(middleware.AuthTokenMiddleware(m.authService))
		{
			protected.GET("", m.handler.List)
			protected.GET("/:id", m.handler.GetByID)
			protected.PATCH("/:id", m.handler.Update)
			protected.DELETE("/:id", m.handler.Delete)
		}
	}
}
