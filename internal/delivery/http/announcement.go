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

// Create creates a new announcement
// @Summary Create announcement
// @Description Create a new course announcement (Teacher/Admin only)
// @Tags announcements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body struct{CourseID string `json:"course_id" binding:"required"`; Title string `json:"title" binding:"required"`; Content string `json:"content" binding:"required"`; Pinned bool `json:"pinned"`} true "Create Announcement Request"
// @Success 201 {object} announcement.Announcement "Created announcement"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/announcements [post]
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

// GetByID returns an announcement by ID
// @Summary Get announcement by ID
// @Description Returns announcement details by ID
// @Tags announcements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Announcement ID"
// @Success 200 {object} announcement.Announcement "Announcement details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Announcement not found"
// @Router /api/v1/announcements/{id} [get]
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

// ListByCourse returns announcements for a specific course
// @Summary List course announcements
// @Description Returns a paginated list of announcements for a course
// @Tags announcements
// @Security BearerAuth
// @Produce json
// @Param courseID path string true "Course ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Announcements list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/courses/{courseID}/announcements [get]
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

// Update updates announcement details
// @Summary Update announcement
// @Description Update announcement details (Author/Admin only)
// @Tags announcements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Announcement ID"
// @Param request body struct{Title *string `json:"title"`; Content *string `json:"content"`; Pinned *bool `json:"pinned"`} true "Update Request"
// @Success 200 {object} announcement.Announcement "Updated announcement"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Announcement not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/announcements/{id} [put]
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

// Delete deletes an announcement
// @Summary Delete announcement
// @Description Deletes an announcement by ID (Author/Admin only)
// @Tags announcements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Announcement ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Announcement not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/announcements/{id} [delete]
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
