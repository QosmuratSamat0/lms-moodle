package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/submission"
	submissionUC "github.com/ap1-final-mini-moodle/internal/usecase/submission"
	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	service *submissionUC.Service
}

func NewSubmissionHandler(service *submissionUC.Service) *SubmissionHandler {
	return &SubmissionHandler{service: service}
}

func (h *SubmissionHandler) Submit(c *gin.Context) {
	var req struct {
		AssignmentID string  `json:"assignment_id" binding:"required"`
		StudentID    string  `json:"student_id" binding:"required"`
		Content      string  `json:"content" binding:"required"`
		FileURL      *string `json:"file_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.service.Submit(&submission.CreateSubmissionInput{
		AssignmentID: req.AssignmentID,
		StudentID:    req.StudentID,
		Content:      req.Content,
		FileURL:      req.FileURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *SubmissionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	s, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *SubmissionHandler) ListByAssignment(c *gin.Context) {
	assignmentID := c.Param("assignmentID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	submissions, err := h.service.ListByAssignment(assignmentID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, submissions)
}

func (h *SubmissionHandler) ListByStudent(c *gin.Context) {
	studentID := c.Param("studentID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	submissions, err := h.service.ListByStudent(studentID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, submissions)
}

func (h *SubmissionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
