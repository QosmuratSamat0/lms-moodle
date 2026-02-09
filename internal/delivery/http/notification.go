package http

import (
	"net/http"

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
// @Param userID path string false "User ID (optional, normally from context)"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} notification.Notification "Notifications list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) ListByUser(c *gin.Context) {
	userID := c.Param("userID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notifications, err := h.service.ListByUser(userID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
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
	c.JSON(http.StatusOK, nil)
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
