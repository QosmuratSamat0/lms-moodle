package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/chat"
	chatUC "github.com/ap1-final-mini-moodle/internal/usecase/chat"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service *chatUC.Service
}

type CreateMessageRequest struct {
	SenderID string `json:"sender_id" binding:"required"`
	CourseID string `json:"course_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func NewChatHandler(service *chatUC.Service) *ChatHandler {
	return &ChatHandler{service: service}
}

// SendMessage sends a new message to a course chat
// @Summary Send message
// @Description Sends a new message to a course chat
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateMessageRequest true "Message Request"
// @Success 201 {object} chat.Message "Created message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/chat [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.service.SendMessage(&chat.CreateMessageInput{
		SenderID: req.SenderID,
		CourseID: req.CourseID,
		Content:  req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

// ListByCourse returns messages for a specific course chat
// @Summary List chat messages
// @Description Returns a paginated list of chat messages for a course
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Param courseID path string true "Course ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(20)
// @Success 200 {array} chat.Message "Messages list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/chat/course/{courseID} [get]
func (h *ChatHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=20"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messages, err := h.service.ListByCourse(courseID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// Delete deletes a chat message
// @Summary Delete message
// @Description Deletes a chat message by ID (Admin only)
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Message not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/chat/{id} [delete]
func (h *ChatHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
