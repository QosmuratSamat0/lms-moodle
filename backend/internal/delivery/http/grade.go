package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/grade"
	gradeUC "github.com/ap1-final-mini-moodle/internal/usecase/grade"
	"github.com/gin-gonic/gin"
)

type GradeHandler struct {
	service *gradeUC.Service
}

func NewGradeHandler(service *gradeUC.Service) *GradeHandler {
	return &GradeHandler{service: service}
}

func (h *GradeHandler) Grade(c *gin.Context) {
	var req struct {
		SubmissionID string `json:"submission_id" binding:"required"`
		Score        int    `json:"score" binding:"required,min=0"`
		Feedback     string `json:"feedback"`
		GradedBy     string `json:"graded_by" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g, err := h.service.Grade(&grade.CreateGradeInput{
		SubmissionID: req.SubmissionID,
		Score:        req.Score,
		Feedback:     req.Feedback,
		GradedBy:     req.GradedBy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *GradeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	g, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *GradeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
