package group

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles group-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new group handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles group creation
// @Summary Create a new group
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateGroupRequest true "Group details"
// @Success 201 {object} GroupResponse
// @Failure 400,409 {object} utils.Response
// @Router /groups [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	group, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, group)
}

// GetByID retrieves a group by ID
// @Summary Get group by ID
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} GroupResponse
// @Failure 404 {object} utils.Response
// @Router /groups/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid group ID")
		return
	}

	group, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, group)
}

// GetByCode retrieves a group by code
// @Summary Get group by code
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param code path string true "Group code (e.g., SE-2430)"
// @Success 200 {object} GroupResponse
// @Failure 404 {object} utils.Response
// @Router /groups/code/{code} [get]
func (h *Handler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		utils.BadRequest(c, "group code is required")
		return
	}

	group, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, group)
}

// Update updates a group
// @Summary Update group
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Param id path string true "Group ID"
// @Param request body UpdateGroupRequest true "Update details"
// @Success 200 {object} utils.Response
// @Failure 400,404,409 {object} utils.Response
// @Router /groups/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid group ID")
		return
	}

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "group updated successfully"})
}

// Delete deletes a group
// @Summary Delete group
// @Tags groups
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 204
// @Failure 404 {object} utils.Response
// @Router /groups/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid group ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// List lists all groups
// @Summary List groups
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} GroupListResponse
// @Router /groups [get]
func (h *Handler) List(c *gin.Context) {
	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.List(c.Request.Context(), pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// AssignTeacher assigns a teacher to a course-group
// @Summary Assign teacher to course-group
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body AssignTeacherRequest true "Assignment details"
// @Success 201 {object} TeacherAssignmentResponse
// @Failure 400,409 {object} utils.Response
// @Router /groups/assignments [post]
func (h *Handler) AssignTeacher(c *gin.Context) {
	var req AssignTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	assignment, err := h.service.AssignTeacher(c.Request.Context(), &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, assignment)
}

// UnassignTeacher removes a teacher assignment
// @Summary Remove teacher assignment
// @Tags groups
// @Security BearerAuth
// @Param id path string true "Assignment ID"
// @Success 204
// @Failure 404 {object} utils.Response
// @Router /groups/assignments/{id} [delete]
func (h *Handler) UnassignTeacher(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid assignment ID")
		return
	}

	if err := h.service.UnassignTeacher(c.Request.Context(), id); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTeacherAssignments lists assignments for the current teacher (from JWT)
// @Summary List my course-group assignments
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} TeacherAssignmentListResponse
// @Router /teachers/me/group-assignments [get]
func (h *Handler) ListTeacherAssignments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	teacherID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalError(c, "invalid user ID format")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListTeacherAssignments(c.Request.Context(), teacherID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// ListTeacherAssignmentsByID lists assignments for a specific teacher (admin only)
// @Summary List teacher's course-group assignments by ID
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} TeacherAssignmentListResponse
// @Router /teachers/{id}/group-assignments [get]
func (h *Handler) ListTeacherAssignmentsByID(c *gin.Context) {
	teacherIDStr := c.Param("id")
	teacherID, err := uuid.Parse(teacherIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid teacher ID")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListTeacherAssignments(c.Request.Context(), teacherID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// ListGroupAssignments lists assignments for a group
// @Summary List group's teacher-course assignments
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} TeacherAssignmentListResponse
// @Router /groups/{id}/assignments [get]
func (h *Handler) ListGroupAssignments(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid group ID")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListGroupAssignments(c.Request.Context(), groupID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}
