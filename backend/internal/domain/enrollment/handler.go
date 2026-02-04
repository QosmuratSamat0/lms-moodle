package enrollment

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles enrollment-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new enrollment handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Enroll handles student enrollment in a course
// @Summary Enroll in a course
// @Tags enrollments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body EnrollRequest true "Enrollment details"
// @Success 201 {object} EnrollmentResponse
// @Failure 400,409 {object} utils.Response
// @Router /enrollments [post]
func (h *Handler) Enroll(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req EnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	enrollment, err := h.service.Enroll(c.Request.Context(), userID, req.CourseID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, enrollment)
}

// GetByID retrieves an enrollment by ID
// @Summary Get enrollment by ID
// @Tags enrollments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 200 {object} EnrollmentResponse
// @Failure 404 {object} utils.Response
// @Router /enrollments/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid enrollment ID")
		return
	}

	enrollment, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, enrollment)
}

// UpdateStatus updates an enrollment status (teacher/admin)
// @Summary Update enrollment status
// @Tags enrollments
// @Security BearerAuth
// @Accept json
// @Param id path string true "Enrollment ID"
// @Param request body UpdateEnrollmentStatusRequest true "Status update"
// @Success 204
// @Failure 400,403,404 {object} utils.Response
// @Router /enrollments/{id}/status [put]
func (h *Handler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid enrollment ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	var req UpdateEnrollmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), id, req.Status, userID, role); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Drop allows a student to drop their enrollment
// @Summary Drop enrollment
// @Tags enrollments
// @Security BearerAuth
// @Param id path string true "Enrollment ID"
// @Success 204
// @Failure 403,404 {object} utils.Response
// @Router /enrollments/{id}/drop [post]
func (h *Handler) Drop(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid enrollment ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	if err := h.service.Drop(c.Request.Context(), id, userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListByCourse lists enrollments for a course
// @Summary List course enrollments
// @Tags enrollments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} EnrollmentListResponse
// @Router /courses/{id}/enrollments [get]
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

// ListMyEnrollments lists enrollments for the current student
// @Summary List my enrollments
// @Tags enrollments
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} EnrollmentListResponse
// @Router /student/enrollments [get]
func (h *Handler) ListMyEnrollments(c *gin.Context) {
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
