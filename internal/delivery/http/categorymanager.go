package http

import (
	"net/http"
	"strconv"

	"github.com/ap1-final-mini-moodle/internal/domain/categorymanager"
	categorymanagerUC "github.com/ap1-final-mini-moodle/internal/usecase/categorymanager"
	"github.com/gin-gonic/gin"
)

type CategoryManagerHandler struct {
	service *categorymanagerUC.Service
}

func NewCategoryManagerHandler(service *categorymanagerUC.Service) *CategoryManagerHandler {
	return &CategoryManagerHandler{service: service}
}

// Create creates a new category manager assignment
// @Summary Create category manager
// @Description Assigns a user as a manager for a specific category (Admin only)
// @Tags category-managers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body categorymanager.CreateCategoryManagerInput true "Create Category Manager Request"
// @Success 201 {object} categorymanager.CategoryManager
// @Router /api/v1/category-managers [post]
func (h *CategoryManagerHandler) Create(c *gin.Context) {
	var req categorymanager.CreateCategoryManagerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cm, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cm)
}

// GetByID returns a category manager by ID
// @Summary Get category manager by ID
// @Description Gets a specific category manager assignment
// @Tags category-managers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category Manager ID"
// @Success 200 {object} categorymanager.CategoryManager
// @Router /api/v1/category-managers/{id} [get]
func (h *CategoryManagerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	cm, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cm)
}

// GetByUserAndCategory returns a category manager assignment for a user and category
// @Summary Get assignment by user and category
// @Description Gets a specific category manager assignment using user and category IDs
// @Tags category-managers
// @Security BearerAuth
// @Produce json
// @Param user_id query string true "User ID"
// @Param category_id query string true "Category ID"
// @Success 200 {object} categorymanager.CategoryManager
// @Router /api/v1/category-managers/user-category [get]
func (h *CategoryManagerHandler) GetByUserAndCategory(c *gin.Context) {
	userID := c.Query("user_id")
	categoryID := c.Query("category_id")

	if userID == "" || categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and category_id are required"})
		return
	}

	cm, err := h.service.GetByUserAndCategory(c.Request.Context(), userID, categoryID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cm)
}

// GetByUserID returns all category manager assignments for a specific user
// @Summary List assignments by user
// @Description Gets all category manager records for a user
// @Tags category-managers
// @Security BearerAuth
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} categorymanager.CategoryManager
// @Router /api/v1/category-managers/by-user/{userID} [get]
func (h *CategoryManagerHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("userID")
	managers, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, managers)
}

// GetByCategoryID returns all category manager assignments for a specific category
// @Summary List assignments by category
// @Description Gets all managers for a specific category
// @Tags category-managers
// @Security BearerAuth
// @Produce json
// @Param categoryID path string true "Category ID"
// @Success 200 {array} categorymanager.CategoryManager
// @Router /api/v1/category-managers/by-category/{categoryID} [get]
func (h *CategoryManagerHandler) GetByCategoryID(c *gin.Context) {
	categoryID := c.Param("categoryID")
	managers, err := h.service.GetByCategoryID(c.Request.Context(), categoryID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, managers)
}

// List returns all category manager assignments with pagination
// @Summary List all assignments
// @Description Gets a paginated list of all category manager assignments
// @Tags category-managers
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Returns data (array of CategoryManager), limit, and offset"
// @Router /api/v1/category-managers [get]
func (h *CategoryManagerHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	managers, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   managers,
		"limit":  limit,
		"offset": offset,
	})
}

// Update updates a category manager assignment
// @Summary Update assignment
// @Description Updates permissions or details of a category manager (Admin/Manager only)
// @Tags category-managers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category Manager ID"
// @Param request body categorymanager.UpdateCategoryManagerInput true "Update Category Manager Request"
// @Success 200 {object} categorymanager.CategoryManager
// @Router /api/v1/category-managers/{id} [put]
func (h *CategoryManagerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req categorymanager.UpdateCategoryManagerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cm, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cm)
}

// Delete removes a category manager assignment
// @Summary Delete assignment
// @Description Removes a category manager assignment (Admin only)
// @Tags category-managers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category Manager ID"
// @Success 204 "No Content"
// @Router /api/v1/category-managers/{id} [delete]
func (h *CategoryManagerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetWithDetails returns a category manager with additional details
// @Summary Get assignment with details
// @Description Gets a specific category manager assignment including related data
// @Tags category-managers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category Manager ID"
// @Success 200 {object} categorymanager.CategoryManager
// @Router /api/v1/category-managers/{id}/details [get]
func (h *CategoryManagerHandler) GetWithDetails(c *gin.Context) {
	id := c.Param("id")
	cm, err := h.service.GetWithDetails(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cm)
}
