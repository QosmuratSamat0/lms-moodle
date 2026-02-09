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

// Grade grades a submission
// @Summary Grade submission
// @Description Grades a student submission (Teacher/Admin only)
// @Tags grades
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body struct{SubmissionID string `json:"submission_id" binding:"required"`; Score int `json:"score" binding:"required,min=0"`; Feedback string `json:"feedback"`; GradedBy string `json:"graded_by" binding:"required"`} true "Grade Request"
// @Success 201 {object} grade.Grade "Created grade"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/grades [post]
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

// GetByID returns a grade by ID
// @Summary Get grade by ID
// @Description Returns grade details by ID
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param id path string true "Grade ID"
// @Success 200 {object} grade.Grade "Grade details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Grade not found"
// @Router /api/v1/grades/{id} [get]
func (h *GradeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	g, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}
	c.JSON(http.StatusOK, g)
}

// Delete deletes a grade
// @Summary Delete grade
// @Description Deletes a grade by ID (Admin only)
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param id path string true "Grade ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Grade not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/grades/{id} [delete]
func (h *GradeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
