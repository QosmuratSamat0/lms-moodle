package assignment

import (
	"net/http"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles assignment-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new assignment handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles assignment creation (teacher only)
// @Summary Create a new assignment
// @Tags assignments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateAssignmentRequest true "Assignment details"
// @Success 201 {object} AssignmentResponse
// @Failure 400,403 {object} utils.Response
// @Router /assignments [post]
func (h *Handler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	assignment, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, assignment)
}

// GetByID retrieves an assignment by ID
// @Summary Get assignment by ID
// @Tags assignments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} AssignmentResponse
// @Failure 404 {object} utils.Response
// @Router /assignments/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid assignment ID")
		return
	}

	assignment, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, assignment)
}

// Update updates an assignment
// @Summary Update assignment
// @Tags assignments
// @Security BearerAuth
// @Accept json
// @Param id path string true "Assignment ID"
// @Param request body UpdateAssignmentRequest true "Update details"
// @Success 200 {object} utils.Response
// @Failure 400,403,404 {object} utils.Response
// @Router /assignments/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid assignment ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	var req UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(c.Request.Context(), id, userID, role, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "assignment updated successfully"})
}

// Delete deletes an assignment
// @Summary Delete assignment
// @Tags assignments
// @Security BearerAuth
// @Param id path string true "Assignment ID"
// @Success 204
// @Failure 403,404 {object} utils.Response
// @Router /assignments/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid assignment ID")
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

// ListByCourse lists assignments for a course
// @Summary List course assignments
// @Tags assignments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} AssignmentListResponse
// @Router /courses/{id}/assignments [get]
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

// ListMyAssignments lists assignments created by the current teacher
// @Summary List my assignments (teacher)
// @Tags assignments
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} AssignmentListResponse
// @Router /teacher/assignments [get]
func (h *Handler) ListMyAssignments(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListByTeacher(c.Request.Context(), userID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// ListUpcoming lists upcoming assignments for the current student
// @Summary List upcoming assignments (student)
// @Tags assignments
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} AssignmentResponse
// @Router /student/assignments/upcoming [get]
func (h *Handler) ListUpcoming(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)
	limit := pagination.Limit
	if limit > 50 {
		limit = 50
	}

	assignments, err := h.service.ListUpcoming(c.Request.Context(), userID, limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, assignments)
}
