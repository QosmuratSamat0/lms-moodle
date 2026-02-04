package student

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles student HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new student handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetProfile godoc
// @Summary Get my student profile
// @Tags students
// @Produce json
// @Success 200 {object} StudentResponse
// @Router /students/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	student, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, student)
}

// GetByID godoc
// @Summary Get student by ID
// @Tags students
// @Produce json
// @Param id path string true "Student User ID"
// @Success 200 {object} StudentResponse
// @Router /students/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid student id"))
		return
	}

	student, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, student)
}

// List godoc
// @Summary List all students
// @Tags students
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param group_name query string false "Filter by group"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name or email"
// @Success 200 {object} StudentListResponse
// @Router /students [get]
func (h *Handler) List(c *gin.Context) {
	page, limit := utils.GetPagination(c)

	var filter StudentFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		errorx.HandleError(c, errorx.NewValidationError(err.Error()))
		return
	}

	result, err := h.service.List(c.Request.Context(), &filter, page, limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, result)
}

// UpdateProfile godoc
// @Summary Update my student profile
// @Tags students
// @Accept json
// @Produce json
// @Param request body UpdateStudentRequest true "Update request"
// @Success 200 {object} map[string]string
// @Router /students/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorx.HandleError(c, errorx.NewValidationError(err.Error()))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	if err := h.service.Update(c.Request.Context(), userID, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

// GetStats godoc
// @Summary Get my statistics
// @Tags students
// @Produce json
// @Success 200 {object} StudentStats
// @Router /students/me/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	userID := middleware.GetUserID(c)

	stats, err := h.service.GetStats(c.Request.Context(), userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, stats)
}

// GetStatsByID godoc
// @Summary Get student statistics by ID
// @Tags students
// @Produce json
// @Param id path string true "Student User ID"
// @Success 200 {object} StudentStats
// @Router /students/{id}/stats [get]
func (h *Handler) GetStatsByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid student id"))
		return
	}

	stats, err := h.service.GetStats(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, stats)
}

// ListByGroup godoc
// @Summary List students by group
// @Tags students
// @Produce json
// @Param group path string true "Group name"
// @Success 200 {array} StudentResponse
// @Router /students/group/{group} [get]
func (h *Handler) ListByGroup(c *gin.Context) {
	groupName := c.Param("group")
	if groupName == "" {
		errorx.HandleError(c, errorx.NewValidationError("group name required"))
		return
	}

	students, err := h.service.ListByGroup(c.Request.Context(), groupName)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, students)
}

// ListByCourse godoc
// @Summary List students enrolled in a course
// @Tags students
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {array} StudentResponse
// @Router /students/course/{id} [get]
func (h *Handler) ListByCourse(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid course id"))
		return
	}

	students, err := h.service.ListByCourse(c.Request.Context(), courseID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, students)
}
