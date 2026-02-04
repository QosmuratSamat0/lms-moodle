package user

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles user-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new user handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Register handles user registration
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} UserResponse
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, user)
}

// Login handles user login
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, resp)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} utils.Response
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, resp)
}

// GetProfile returns the current user's profile
// @Summary Get current user profile
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} utils.Response
// @Router /users/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		utils.Unauthorized(c, "invalid user context")
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, profile)
}

// UpdateProfile updates the current user's profile
// @Summary Update current user profile
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "Profile update"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /users/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		utils.Unauthorized(c, "invalid user context")
		return
	}

	role, _ := c.Get(string(middleware.RoleKey))
	roleStr, _ := role.(string)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.UpdateProfile(c.Request.Context(), userID, roleStr, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "profile updated successfully"})
}

// ChangePassword changes the current user's password
// @Summary Change password
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Password change"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /users/me/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		utils.Unauthorized(c, "invalid user context")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "password changed successfully"})
}

// GetUser returns a user by ID (admin only)
// @Summary Get user by ID
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404 {object} utils.Response
// @Router /admin/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid user ID")
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, user)
}

// ListUsers returns a paginated list of users (admin only)
// @Summary List all users
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} UserListResponse
// @Router /admin/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	pagination := utils.GetPaginationFromContext(c)

	resp, err := h.service.List(c.Request.Context(), pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, resp)
}

// DeactivateUser deactivates a user (admin only)
// @Summary Deactivate user
// @Tags admin
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 404 {object} utils.Response
// @Router /admin/users/{id}/deactivate [post]
func (h *Handler) DeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid user ID")
		return
	}

	if err := h.service.Deactivate(c.Request.Context(), id); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ActivateUser activates a user (admin only)
// @Summary Activate user
// @Tags admin
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 404 {object} utils.Response
// @Router /admin/users/{id}/activate [post]
func (h *Handler) ActivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid user ID")
		return
	}

	if err := h.service.Activate(c.Request.Context(), id); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Helper to extract user ID from Gin context
func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDVal, exists := c.Get(string(middleware.UserIDKey))
	if !exists {
		return uuid.Nil, errorx.NewUnauthorizedError("user ID not found in context")
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return uuid.Nil, errorx.NewUnauthorizedError("invalid user ID type")
	}

	return userID, nil
}
