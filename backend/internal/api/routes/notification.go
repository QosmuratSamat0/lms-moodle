// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/notification"
	"github.com/gin-gonic/gin"
)

// NotificationRouter handles notification-related routes
type NotificationRouter struct {
	notificationHandler *notification.Handler
}

// NewNotificationRouter creates a new notification router
func NewNotificationRouter(notificationHandler *notification.Handler) *NotificationRouter {
	return &NotificationRouter{
		notificationHandler: notificationHandler,
	}
}

// SetupRoutes configures notification routes
func (nr *NotificationRouter) SetupRoutes(api *gin.RouterGroup) {
	notifications := api.Group("/notifications")
	{
		notifications.GET("", nr.notificationHandler.ListNotifications)
		notifications.GET("/unread-count", nr.notificationHandler.GetUnreadCount)
		notifications.GET("/:id", nr.notificationHandler.GetNotification)
		notifications.POST("/:id/read", nr.notificationHandler.MarkAsRead)
		notifications.POST("/read-all", nr.notificationHandler.MarkAllAsRead)
		notifications.DELETE("/:id", nr.notificationHandler.DeleteNotification)
		notifications.DELETE("", nr.notificationHandler.DeleteAllNotifications)
	}
}
