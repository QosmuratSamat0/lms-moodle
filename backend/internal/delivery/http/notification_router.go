package http

import (
	"github.com/gin-gonic/gin"
)

type NotificationModule struct {
	handler *NotificationHandler
	secret  []byte
}

func NewNotificationModule(handler *NotificationHandler, secret []byte) *NotificationModule {
	return &NotificationModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *NotificationModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	notifications := api.Group("/notifications")
	{
		notifications.GET("", m.handler.ListByUser)
		notifications.PATCH("/:id/read", m.handler.MarkAsRead)
		notifications.DELETE("/:id", m.handler.Delete)
	}
}
