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

func NewChatHandler(service *chatUC.Service) *ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req struct {
		SenderID string `json:"sender_id" binding:"required"`
		CourseID string `json:"course_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}
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

func (h *ChatHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
