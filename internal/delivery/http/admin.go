package http

import (
	"net/http"
	"strconv"

	"github.com/ap1-final-mini-moodle/internal/domain/admin"
	adminUC "github.com/ap1-final-mini-moodle/internal/usecase/admin"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *adminUC.Service
}

func NewAdminHandler(service *adminUC.Service) *AdminHandler {
	return &AdminHandler{service: service}
}

// Create creates a new admin profile
// @Summary Create admin
// @Description Create a new admin profile (Super Admin only)
// @Tags admins
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body admin.CreateAdminInput true "Create Admin Request"
// @Success 201 {object} admin.Admin "Created admin"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/admins [post]
func (h *AdminHandler) Create(c *gin.Context) {
	var req admin.CreateAdminInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.service.CreateAdmin(c.Request.Context(), &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

// GetByID returns an admin by ID
// @Summary Get admin by ID
// @Description Returns admin profile by ID (Admin only)
// @Tags admins
// @Security BearerAuth
// @Produce json
// @Param id path string true "Admin ID"
// @Success 200 {object} admin.Admin "Admin details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Admin not found"
// @Router /api/v1/admins/{id} [get]
func (h *AdminHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

// GetByUserID returns an admin by User ID
// @Summary Get admin by User ID
// @Description Returns admin profile by user ID (Admin only)
// @Tags admins
// @Security BearerAuth
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} admin.Admin "Admin details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Admin not found"
// @Router /api/v1/admins/user/{userID} [get]
func (h *AdminHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("userID")
	a, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

// GetMyProfile returns current admin's profile
// @Summary Get my profile
// @Description Returns profile of currently logged-in admin
// @Tags admins
// @Security BearerAuth
// @Produce json
// @Success 200 {object} admin.Admin "Admin details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Admin not found"
// @Router /api/v1/admins/me [get]
func (h *AdminHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	a, err := h.service.GetByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

// List returns a list of admins
// @Summary List admins
// @Description Returns a paginated list of all admins (Admin only)
// @Tags admins
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Admins list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/admins [get]
func (h *AdminHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	admins, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   admins,
		"limit":  limit,
		"offset": offset,
	})
}

// Update updates admin details
// @Summary Update admin
// @Description Update admin profile details (Admin only)
// @Tags admins
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Admin ID"
// @Param request body admin.UpdateAdminInput true "Update Request"
// @Success 200 {object} admin.Admin "Updated admin"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Admin not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/admins/{id} [put]
func (h *AdminHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req admin.UpdateAdminInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.service.UpdateAdmin(c.Request.Context(), id, &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

// Delete deletes an admin profile
// @Summary Delete admin
// @Description Deletes an admin profile by ID (Super Admin only)
// @Tags admins
// @Security BearerAuth
// @Produce json
// @Param id path string true "Admin ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Admin not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/admins/{id} [delete]
func (h *AdminHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteAdmin(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
