package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/announcement"
	announcementUC "github.com/ap1-final-mini-moodle/internal/usecase/announcement"
	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	service *announcementUC.Service
}

func NewAnnouncementHandler(service *announcementUC.Service) *AnnouncementHandler {
	return &AnnouncementHandler{service: service}
}

func (h *AnnouncementHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	authorID := ""
	if userID != nil {
		authorID = userID.(string)
	}

	var req struct {
		CourseID string `json:"course_id" binding:"required"`
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Pinned   bool   `json:"pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.service.Create(c.Request.Context(), authorID, &announcement.CreateAnnouncementInput{
		CourseID: req.CourseID,
		Title:    req.Title,
		Content:  req.Content,
		Pinned:   req.Pinned,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *AnnouncementHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *AnnouncementHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseID")
	var req struct {
		Limit  int `form:"limit,default=20"`
		Offset int `form:"offset,default=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcements, total, err := h.service.GetByCourseID(c.Request.Context(), courseID, req.Limit, req.Offset)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  announcements,
		"total": total,
	})
}

func (h *AnnouncementHandler) Update(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	authorID := ""
	if userID != nil {
		authorID = userID.(string)
	}

	var req struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Pinned  *bool   `json:"pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.service.Update(c.Request.Context(), id, authorID, &announcement.UpdateAnnouncementInput{
		Title:   req.Title,
		Content: req.Content,
		Pinned:  req.Pinned,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	authorID := ""
	if userID != nil {
		authorID = userID.(string)
	}

	if err := h.service.Delete(c.Request.Context(), id, authorID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
