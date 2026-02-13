package http

import (
	"net/http"
	"strconv"

	"github.com/ap1-final-mini-moodle/internal/domain/coursecategory"
	coursecategoryUC "github.com/ap1-final-mini-moodle/internal/usecase/coursecategory"
	"github.com/gin-gonic/gin"
)

type CourseCategoryHandler struct {
	service *coursecategoryUC.Service
}

func NewCourseCategoryHandler(service *coursecategoryUC.Service) *CourseCategoryHandler {
	return &CourseCategoryHandler{service: service}
}

// Create creates a new course category
// @Summary Create course category
// @Description Creates a new course category (Admin only)
// @Tags course-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body coursecategory.CreateCourseCategoryInput true "Create Category Request"
// @Success 201 {object} coursecategory.CourseCategory
// @Router /api/v1/categories [post]
func (h *CourseCategoryHandler) Create(c *gin.Context) {
	var req coursecategory.CreateCourseCategoryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cc, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cc)
}

// GetByID returns a course category by ID
// @Summary Get category by ID
// @Description Gets a specific course category
// @Tags course-categories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} coursecategory.CourseCategory
// @Router /api/v1/categories/{id} [get]
func (h *CourseCategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	cc, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cc)
}

// List returns all course categories with pagination
// @Summary List all categories
// @Description Gets a paginated list of all course categories
// @Tags course-categories
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Returns data (array of CourseCategory), limit, and offset"
// @Router /api/v1/categories [get]
func (h *CourseCategoryHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	categories, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   categories,
		"limit":  limit,
		"offset": offset,
	})
}

// Update updates a course category
// @Summary Update category
// @Description Updates an existing course category (Admin only)
// @Tags course-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body coursecategory.UpdateCourseCategoryInput true "Update Category Request"
// @Success 200 {object} coursecategory.CourseCategory
// @Router /api/v1/categories/{id} [put]
func (h *CourseCategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req coursecategory.UpdateCourseCategoryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cc, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cc)
}

// Delete removes a course category
// @Summary Delete category
// @Description Removes a course category (Admin only)
// @Tags course-categories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Router /api/v1/categories/{id} [delete]
func (h *CourseCategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
