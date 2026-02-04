package manager

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles manager HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new manager handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetProfile godoc
// @Summary Get my manager profile
// @Tags managers
// @Produce json
// @Success 200 {object} ManagerResponse
// @Router /managers/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	manager, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, manager)
}

// GetByID godoc
// @Summary Get manager by ID
// @Tags managers
// @Produce json
// @Param id path string true "Manager User ID"
// @Success 200 {object} ManagerResponse
// @Router /managers/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid manager id"))
		return
	}

	manager, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, manager)
}

// List godoc
// @Summary List all managers
// @Tags managers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name or email"
// @Success 200 {object} ManagerListResponse
// @Router /managers [get]
func (h *Handler) List(c *gin.Context) {
	page, limit := utils.GetPagination(c)

	var filter ManagerFilter
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
// @Summary Update my manager profile
// @Tags managers
// @Accept json
// @Produce json
// @Param request body UpdateManagerRequest true "Update request"
// @Success 200 {object} map[string]string
// @Router /managers/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateManagerRequest
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

// GetSystemOverview godoc
// @Summary Get system overview statistics
// @Tags managers
// @Produce json
// @Success 200 {object} SystemOverview
// @Router /managers/overview [get]
func (h *Handler) GetSystemOverview(c *gin.Context) {
	overview, err := h.service.GetSystemOverview(c.Request.Context())
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, overview)
}

// ActivateUser godoc
// @Summary Activate a user account
// @Tags managers
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Router /managers/users/{id}/activate [post]
func (h *Handler) ActivateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid user id"))
		return
	}

	if err := h.service.ActivateUser(c.Request.Context(), id); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user activated"})
}

// DeactivateUser godoc
// @Summary Deactivate a user account
// @Tags managers
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Router /managers/users/{id}/deactivate [post]
func (h *Handler) DeactivateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorx.HandleError(c, errorx.NewValidationError("invalid user id"))
		return
	}

	if err := h.service.DeactivateUser(c.Request.Context(), id); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deactivated"})
}

// BulkUserAction godoc
// @Summary Perform bulk action on users
// @Tags managers
// @Accept json
// @Produce json
// @Param request body BulkUserActionRequest true "Bulk action request"
// @Success 200 {object} map[string]string
// @Router /managers/users/bulk [post]
func (h *Handler) BulkUserAction(c *gin.Context) {
	var req BulkUserActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorx.HandleError(c, errorx.NewValidationError(err.Error()))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	if err := h.service.BulkUserAction(c.Request.Context(), &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bulk action completed"})
}
