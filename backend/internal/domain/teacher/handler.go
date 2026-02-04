package teacher

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles teacher HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new teacher handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetProfile godoc
// @Summary Get my teacher profile
// @Tags teachers
// @Produce json
// @Success 200 {object} TeacherResponse
// @Router /teachers/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	teacher, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, teacher)
}

// GetByID godoc
// @Summary Get teacher by ID
// @Tags teachers
// @Produce json
// @Param id path string true "Teacher User ID"
// @Success 200 {object} TeacherResponse
// @Router /teachers/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid teacher id"))
		return
	}

	teacher, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, teacher)
}

// List godoc
// @Summary List all teachers
// @Tags teachers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param department query string false "Filter by department"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name or email"
// @Success 200 {object} TeacherListResponse
// @Router /teachers [get]
func (h *Handler) List(c *gin.Context) {
	page, limit := utils.GetPagination(c)

	var filter TeacherFilter
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
// @Summary Update my teacher profile
// @Tags teachers
// @Accept json
// @Produce json
// @Param request body UpdateTeacherRequest true "Update request"
// @Success 200 {object} map[string]string
// @Router /teachers/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateTeacherRequest
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
// @Tags teachers
// @Produce json
// @Success 200 {object} TeacherStats
// @Router /teachers/me/stats [get]
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
// @Summary Get teacher statistics by ID
// @Tags teachers
// @Produce json
// @Param id path string true "Teacher User ID"
// @Success 200 {object} TeacherStats
// @Router /teachers/{id}/stats [get]
func (h *Handler) GetStatsByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid teacher id"))
		return
	}

	stats, err := h.service.GetStats(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, stats)
}

// ListByDepartment godoc
// @Summary List teachers by department
// @Tags teachers
// @Produce json
// @Param department path string true "Department name"
// @Success 200 {array} TeacherResponse
// @Router /teachers/department/{department} [get]
func (h *Handler) ListByDepartment(c *gin.Context) {
	department := c.Param("department")
	if department == "" {
		errorx.HandleError(c, errorx.NewValidationError("department name required"))
		return
	}

	teachers, err := h.service.ListByDepartment(c.Request.Context(), department)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, teachers)
}
