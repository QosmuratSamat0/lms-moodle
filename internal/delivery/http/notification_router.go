package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type NotificationModule struct {
	handler     *NotificationHandler
	authService *authUC.Service
}

func NewNotificationModule(handler *NotificationHandler, authService *authUC.Service) *NotificationModule {
	return &NotificationModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *NotificationModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	notifications := api.Group("/notifications")
	notifications.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		notifications.GET("", m.handler.ListByUser)
		notifications.PATCH("/:id/read", m.handler.MarkAsRead)
		notifications.DELETE("/:id", m.handler.Delete)
	}
}
