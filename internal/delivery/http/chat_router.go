package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type ChatModule struct {
	handler     *ChatHandler
	authService *authUC.Service
}

func NewChatModule(handler *ChatHandler, authService *authUC.Service) *ChatModule {
	return &ChatModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *ChatModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	chat := api.Group("/chat")
	chat.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		chat.POST("", m.handler.SendMessage)
		chat.GET("/course/:courseID", m.handler.ListByCourse)
		chat.DELETE("/:id", m.handler.Delete)
	}
}
