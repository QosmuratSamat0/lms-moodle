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

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.MarkAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
