package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/grade"
	gradeUC "github.com/ap1-final-mini-moodle/internal/usecase/grade"
	notificationUC "github.com/ap1-final-mini-moodle/internal/usecase/notification"
	submissionUC "github.com/ap1-final-mini-moodle/internal/usecase/submission"
	"github.com/gin-gonic/gin"
)

type GradeHandler struct {
	service     *gradeUC.Service
	notifSvc    *notificationUC.Service
	submitSvc   *submissionUC.Service
}

type CreateGradeRequest struct {
	SubmissionID string  `json:"submission_id" binding:"required"`
	Score        float64 `json:"score" binding:"min=0"`
	Feedback     string  `json:"feedback"`
	GradedBy     string  `json:"graded_by"`
}

func NewGradeHandler(service *gradeUC.Service, notifSvc *notificationUC.Service, submitSvc *submissionUC.Service) *GradeHandler {
	return &GradeHandler{service: service, notifSvc: notifSvc, submitSvc: submitSvc}
}

// Grade grades a submission
// @Summary Grade submission
// @Description Grades a student submission (Teacher/Admin only)
// @Tags grades
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateGradeRequest true "Grade Request"
// @Success 201 {object} grade.Grade "Created grade"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/grades [post]
func (h *GradeHandler) Grade(c *gin.Context) {
	var req CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-populate graded_by from JWT if not provided
	gradedBy := req.GradedBy
	if gradedBy == "" {
		if uid, exists := c.Get("userID"); exists {
			gradedBy = uid.(string)
		}
	}

	g, err := h.service.Grade(&grade.CreateGradeInput{
		SubmissionID: req.SubmissionID,
		Score:        req.Score,
		Feedback:     req.Feedback,
		GradedBy:     gradedBy,
	})
	if err != nil {
		log.Printf("[GRADE] Create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notify the student that their submission was graded
	if h.notifSvc != nil && h.submitSvc != nil {
		go func() {
			sub, err := h.submitSvc.GetByID(req.SubmissionID)
			if err != nil {
				log.Printf("[NOTIFICATION] Failed to get submission: %v", err)
				return
			}
			title := "Assignment Graded"
			msg := fmt.Sprintf("Your submission has been graded: %.0f points.", req.Score)
			if req.Feedback != "" {
				msg += fmt.Sprintf(" Feedback: %s", req.Feedback)
			}
			h.notifSvc.NotifyMany([]string{sub.StudentID}, "success", title, msg)
		}()
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
