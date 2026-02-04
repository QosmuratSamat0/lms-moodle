package http

import (
	"github.com/gin-gonic/gin"
)

type ChatModule struct {
	handler *ChatHandler
	secret  []byte
}

func NewChatModule(handler *ChatHandler, secret []byte) *ChatModule {
	return &ChatModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *ChatModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	chat := api.Group("/chat")
	{
		chat.POST("", m.handler.SendMessage)
		chat.GET("/course/:courseID", m.handler.ListByCourse)
		chat.DELETE("/:id", m.handler.Delete)
	}
}
