package grade

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles grade-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new grade handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GradeSubmission grades a submission (teacher only)
// @Summary Grade a submission
// @Tags grades
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body GradeSubmissionRequest true "Grade details"
// @Success 201 {object} GradeResponse
// @Failure 400 {object} utils.Response
// @Router /grades [post]
func (h *Handler) GradeSubmission(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	grade, err := h.service.GradeSubmission(c.Request.Context(), userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, grade)
}

// GetByID retrieves a grade by ID
// @Summary Get grade by ID
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param id path string true "Grade ID"
// @Success 200 {object} GradeResponse
// @Failure 404 {object} utils.Response
// @Router /grades/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid grade ID")
		return
	}

	grade, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, grade)
}

// Update updates a grade
// @Summary Update grade
// @Tags grades
// @Security BearerAuth
// @Accept json
// @Param id path string true "Grade ID"
// @Param request body UpdateGradeRequest true "Update details"
// @Success 200 {object} utils.Response
// @Failure 400,403,404 {object} utils.Response
// @Router /grades/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid grade ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	var req UpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(c.Request.Context(), id, userID, role, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "grade updated successfully"})
}

// Delete deletes a grade
// @Summary Delete grade
// @Tags grades
// @Security BearerAuth
// @Param id path string true "Grade ID"
// @Success 204
// @Failure 403,404 {object} utils.Response
// @Router /grades/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid grade ID")
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

// ListMyGrades lists grades for the current student
// @Summary List my grades
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} GradeListResponse
// @Router /student/grades [get]
func (h *Handler) ListMyGrades(c *gin.Context) {
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

// ListByCourse lists grades for a course (teacher/admin)
// @Summary List course grades
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} GradeListResponse
// @Router /courses/{id}/grades [get]
func (h *Handler) ListByCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListByCourse(c.Request.Context(), courseID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// GetStudentCourseSummary gets grade summary for a student in a course
// @Summary Get student course grade summary
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} StudentGradeSummary
// @Router /student/courses/{course_id}/summary [get]
func (h *Handler) GetStudentCourseSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	courseIDStr := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	summary, err := h.service.GetStudentCourseSummary(c.Request.Context(), userID, courseID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, summary)
}
