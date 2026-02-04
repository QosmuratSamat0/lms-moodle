package notification

import (
	"net/http"
	"strconv"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles notification-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new notification handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ListNotifications lists user's notifications
// @Summary List notifications
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Param unread_only query bool false "Show only unread"
// @Success 200 {object} NotificationListResponse
// @Router /notifications [get]
func (h *Handler) ListNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)
	unreadOnly, _ := strconv.ParseBool(c.Query("unread_only"))

	list, err := h.service.List(c.Request.Context(), userID, unreadOnly, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// GetNotification gets a notification by ID
// @Summary Get notification
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} NotificationResponse
// @Router /notifications/{id} [get]
func (h *Handler) GetNotification(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid notification ID")
		return
	}

	notification, err := h.service.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, notification)
}

// MarkAsRead marks a notification as read
// @Summary Mark notification as read
// @Tags notifications
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} utils.Response
// @Router /notifications/{id}/read [post]
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid notification ID")
		return
	}

	if err := h.service.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "notification marked as read"})
}

// MarkAllAsRead marks all notifications as read
// @Summary Mark all notifications as read
// @Tags notifications
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /notifications/read-all [post]
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	if err := h.service.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "all notifications marked as read"})
}

// DeleteNotification deletes a notification
// @Summary Delete notification
// @Tags notifications
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 204
// @Router /notifications/{id} [delete]
func (h *Handler) DeleteNotification(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid notification ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteAllNotifications deletes all notifications
// @Summary Delete all notifications
// @Tags notifications
// @Security BearerAuth
// @Success 204
// @Router /notifications [delete]
func (h *Handler) DeleteAllNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	if err := h.service.DeleteAll(c.Request.Context(), userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUnreadCount gets the count of unread notifications
// @Summary Get unread count
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]int64
// @Router /notifications/unread-count [get]
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	count, err := h.service.CountUnread(c.Request.Context(), userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"unread_count": count})
}
