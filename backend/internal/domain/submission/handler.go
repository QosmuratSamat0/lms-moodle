package submission

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles submission-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new submission handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Submit handles assignment submission
// @Summary Submit an assignment
// @Tags submissions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateSubmissionRequest true "Submission details"
// @Success 201 {object} SubmissionResponse
// @Failure 400,409 {object} utils.Response
// @Router /submissions [post]
func (h *Handler) Submit(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	submission, err := h.service.Submit(c.Request.Context(), userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, submission)
}

// GetByID retrieves a submission by ID
// @Summary Get submission by ID
// @Tags submissions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Submission ID"
// @Success 200 {object} SubmissionResponse
// @Failure 404 {object} utils.Response
// @Router /submissions/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid submission ID")
		return
	}

	submission, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, submission)
}

// Update updates a submission
// @Summary Update submission
// @Tags submissions
// @Security BearerAuth
// @Accept json
// @Param id path string true "Submission ID"
// @Param request body UpdateSubmissionRequest true "Update details"
// @Success 200 {object} utils.Response
// @Failure 400,403,404 {object} utils.Response
// @Router /submissions/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid submission ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req UpdateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(c.Request.Context(), id, userID, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "submission updated successfully"})
}

// Delete deletes a submission
// @Summary Delete submission
// @Tags submissions
// @Security BearerAuth
// @Param id path string true "Submission ID"
// @Success 204
// @Failure 403,404 {object} utils.Response
// @Router /submissions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid submission ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	if err := h.service.Delete(c.Request.Context(), id, userID, role); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListByAssignment lists submissions for an assignment
// @Summary List assignment submissions
// @Tags submissions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Assignment ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Param group_id query string false "Filter by student group ID"
// @Success 200 {object} SubmissionListResponse
// @Router /assignments/{id}/submissions [get]
func (h *Handler) ListByAssignment(c *gin.Context) {
	assignmentIDStr := c.Param("id")
	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid assignment ID")
		return
	}

	// Optional group_id filter
	var groupID *uuid.UUID
	if groupIDStr := c.Query("group_id"); groupIDStr != "" {
		parsedID, err := uuid.Parse(groupIDStr)
		if err != nil {
			utils.BadRequest(c, "invalid group ID")
			return
		}
		groupID = &parsedID
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListByAssignment(c.Request.Context(), assignmentID, groupID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// ListMySubmissions lists submissions for the current student
// @Summary List my submissions
// @Tags submissions
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} SubmissionListResponse
// @Router /student/submissions [get]
func (h *Handler) ListMySubmissions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListByStudent(c.Request.Context(), userID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}
