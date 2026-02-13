package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/notification"
	notificationUC "github.com/ap1-final-mini-moodle/internal/usecase/notification"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *notificationUC.Service
}

func NewNotificationHandler(service *notificationUC.Service) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// ListByUser returns notifications for a specific user
// @Summary List user notifications
// @Description Returns a paginated list of notifications for the current user
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} notification.Notification "Notifications list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) ListByUser(c *gin.Context) {
	// Get userID from JWT claims
	claims, err := GetTokenClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=50"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notifications, err := h.service.ListByUser(claims.UserID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if notifications == nil {
		notifications = []*notification.Notification{}
	}

	// Count unread
	unreadCount := 0
	for _, n := range notifications {
		if !n.Read {
			unreadCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         len(notifications),
		"unread_count":  unreadCount,
	})
}

// MarkAsRead marks a notification as read
// @Summary Mark notification as read
// @Description Updates notification status to read
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Notification not found"
// @Router /api/v1/notifications/{id}/read [patch]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.MarkAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	claims, err := GetTokenClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.MarkAllAsRead(claims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all marked as read"})
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	claims, err := GetTokenClaimsFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	count, err := h.service.UnreadCount(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// Delete deletes a notification
// @Summary Delete notification
// @Description Deletes a notification by ID
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param id path string true "Notification ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Notification not found"
// @Router /api/v1/notifications/{id} [delete]
func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
