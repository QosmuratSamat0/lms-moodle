package plagiarism

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles plagiarism HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new plagiarism handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CheckSubmissionRequest represents a request to check submission
type CheckSubmissionRequest struct {
	SubmissionID uuid.UUID `json:"submission_id" validate:"required"`
	Text         string    `json:"text" validate:"required"`
}

// CheckSubmission godoc
// @Summary Run plagiarism check on a submission
// @Tags plagiarism
// @Accept json
// @Produce json
// @Param request body CheckSubmissionRequest true "Check request"
// @Success 200 {object} Report
// @Router /plagiarism/check [post]
func (h *Handler) CheckSubmission(c *gin.Context) {
	var req CheckSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorx.HandleError(c, errorx.NewValidationError(err.Error()))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	report, err := h.service.CheckSubmission(c.Request.Context(), req.SubmissionID, req.Text)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, report)
}

// GetReport godoc
// @Summary Get plagiarism report by ID
// @Tags plagiarism
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} ReportWithSubmission
// @Router /plagiarism/reports/{id} [get]
func (h *Handler) GetReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid report id"))
		return
	}

	report, err := h.service.GetReport(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, report)
}

// GetReportBySubmission godoc
// @Summary Get plagiarism report by submission ID
// @Tags plagiarism
// @Produce json
// @Param submission_id path string true "Submission ID"
// @Success 200 {object} Report
// @Router /plagiarism/submissions/{submission_id} [get]
func (h *Handler) GetReportBySubmission(c *gin.Context) {
	submissionID, err := uuid.Parse(c.Param("submission_id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid submission id"))
		return
	}

	report, err := h.service.GetReportBySubmission(c.Request.Context(), submissionID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, report)
}

// ListByAssignment godoc
// @Summary List plagiarism reports for an assignment
// @Tags plagiarism
// @Produce json
// @Param assignment_id path string true "Assignment ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} ReportListResponse
// @Router /plagiarism/assignments/{assignment_id} [get]
func (h *Handler) ListByAssignment(c *gin.Context) {
	assignmentID, err := uuid.Parse(c.Param("assignment_id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid assignment id"))
		return
	}

	page, limit := utils.GetPagination(c)

	result, err := h.service.ListByAssignment(c.Request.Context(), assignmentID, page, limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, result)
}

// ListSuspicious godoc
// @Summary List suspicious submissions above threshold
// @Tags plagiarism
// @Produce json
// @Param threshold query number false "Similarity threshold" default(30)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} ReportListResponse
// @Router /plagiarism/suspicious [get]
func (h *Handler) ListSuspicious(c *gin.Context) {
	var threshold *float64
	if t := c.Query("threshold"); t != "" {
		var parsed float64
		if _, err := utils.ParseFloat(t, &parsed); err == nil {
			threshold = &parsed
		}
	}

	page, limit := utils.GetPagination(c)

	result, err := h.service.ListSuspicious(c.Request.Context(), threshold, page, limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, result)
}

// Recheck godoc
// @Summary Recheck a submission for plagiarism
// @Tags plagiarism
// @Param submission_id path string true "Submission ID"
// @Success 200 {object} Report
// @Router /plagiarism/submissions/{submission_id}/recheck [post]
func (h *Handler) Recheck(c *gin.Context) {
	submissionID, err := uuid.Parse(c.Param("submission_id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid submission id"))
		return
	}

	report, err := h.service.Recheck(c.Request.Context(), submissionID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report, "message": "recheck completed"})
}
