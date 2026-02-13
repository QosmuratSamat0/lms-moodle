package http

import (
	"net/http"
	"strconv"

	"github.com/ap1-final-mini-moodle/internal/domain/manager"
	managerUC "github.com/ap1-final-mini-moodle/internal/usecase/manager"
	"github.com/gin-gonic/gin"
)

type ManagerHandler struct {
	service *managerUC.Service
}

type AddManagedCategoryRequest struct {
	CategoryID string `json:"category_id" binding:"required"`
}

type AddManagedTeacherRequest struct {
	TeacherID string `json:"teacher_id" binding:"required"`
}

type UpdateManagedResourcesRequest struct {
	Categories []string `json:"categories,omitempty"`
	Teachers   []string `json:"teachers,omitempty"`
}

func NewManagerHandler(service *managerUC.Service) *ManagerHandler {
	return &ManagerHandler{service: service}
}

// Create creates a new manager
// @Summary Create manager
// @Description Creates a new manager profile (Admin only)
// @Tags managers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body manager.CreateManagerInput true "Create Manager Request"
// @Success 201 {object} manager.Manager
// @Router /api/v1/managers [post]
func (h *ManagerHandler) Create(c *gin.Context) {
	var req manager.CreateManagerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.service.CreateManager(c.Request.Context(), &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

// GetByID returns a manager by ID
// @Summary Get manager by ID
// @Description Gets a specific manager profile
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Manager ID"
// @Success 200 {object} manager.Manager
// @Router /api/v1/managers/{id} [get]
func (h *ManagerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	m, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// GetByUserID returns a manager by user ID
// @Summary Get manager by user ID
// @Description Gets a manager profile associated with a user
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} manager.Manager
// @Router /api/v1/managers/user/{userID} [get]
func (h *ManagerHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("userID")
	m, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// GetMyProfile returns the profile of the authenticated manager
// @Summary Get my profile
// @Description Gets the current manager's profile
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} manager.Manager
// @Router /api/v1/managers/me [get]
func (h *ManagerHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	m, err := h.service.GetByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// GetByDepartment returns managers by department
// @Summary List managers by department
// @Description Gets all managers in a specific department
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Param department query string true "Department Name"
// @Success 200 {array} manager.Manager
// @Router /api/v1/managers/department [get]
func (h *ManagerHandler) GetByDepartment(c *gin.Context) {
	department := c.Query("department")
	if department == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "department is required"})
		return
	}

	managers, err := h.service.GetByDepartment(c.Request.Context(), department)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, managers)
}

// List returns all managers with pagination
// @Summary List all managers
// @Description Gets a paginated list of all manager profiles
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Returns data (array of Manager), limit, and offset"
// @Router /api/v1/managers [get]
func (h *ManagerHandler) List(c *gin.Context) {
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

// Update updates a manager profile
// @Summary Update manager
// @Description Updates an existing manager profile (Admin only)
// @Tags managers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Manager ID"
// @Param request body manager.UpdateManagerInput true "Update Manager Request"
// @Success 200 {object} manager.Manager
// @Router /api/v1/managers/{id} [put]
func (h *ManagerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req manager.UpdateManagerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.service.UpdateManager(c.Request.Context(), id, &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// Delete removes a manager profile
// @Summary Delete manager
// @Description Removes a manager profile (Admin only)
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Manager ID"
// @Success 204 "No Content"
// @Router /api/v1/managers/{id} [delete]
func (h *ManagerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteManager(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetWithDetails returns a manager with related entities
// @Summary Get manager with details
// @Description Gets a specific manager profile including managed categories and teachers
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Manager ID"
// @Success 200 {object} manager.Manager
// @Router /api/v1/managers/{id}/details [get]
func (h *ManagerHandler) GetWithDetails(c *gin.Context) {
	id := c.Param("id")
	m, err := h.service.GetWithDetails(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// AddManagedCategory assigns a category to a manager
// @Summary Add managed category
// @Description Assigns a course category to be managed by this manager (Admin only)
// @Tags managers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Manager ID"
// @Param request body AddManagedCategoryRequest true "Add Category Request"
// @Success 204 "No Content"
// @Router /api/v1/managers/{id}/categories [post]
func (h *ManagerHandler) AddManagedCategory(c *gin.Context) {
	managerID := c.Param("id")
	var req AddManagedCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddManagedCategory(c.Request.Context(), managerID, req.CategoryID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// RemoveManagedCategory unassigns a category from a manager
// @Summary Remove managed category
// @Description Removes a course category from this manager's responsibilities (Admin only)
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Manager ID"
// @Param categoryId path string true "Category ID"
// @Success 204 "No Content"
// @Router /api/v1/managers/{id}/categories/{categoryId} [delete]
func (h *ManagerHandler) RemoveManagedCategory(c *gin.Context) {
	managerID := c.Param("id")
	categoryID := c.Param("categoryId")

	if err := h.service.RemoveManagedCategory(c.Request.Context(), managerID, categoryID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// AddManagedTeacher assigns a teacher to a manager
// @Summary Add managed teacher
// @Description Assigns a teacher to be managed by this manager (Admin only)
// @Tags managers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Manager ID"
// @Param request body AddManagedTeacherRequest true "Add Teacher Request"
// @Success 204 "No Content"
// @Router /api/v1/managers/{id}/teachers [post]
func (h *ManagerHandler) AddManagedTeacher(c *gin.Context) {
	managerID := c.Param("id")
	var req AddManagedTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddManagedTeacher(c.Request.Context(), managerID, req.TeacherID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// RemoveManagedTeacher unassigns a teacher from a manager
// @Summary Remove managed teacher
// @Description Removes a teacher from this manager's responsibilities (Admin only)
// @Tags managers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Manager ID"
// @Param teacherId path string true "Teacher ID"
// @Success 204 "No Content"
// @Router /api/v1/managers/{id}/teachers/{teacherId} [delete]
func (h *ManagerHandler) RemoveManagedTeacher(c *gin.Context) {
	managerID := c.Param("id")
	teacherID := c.Param("teacherId")

	if err := h.service.RemoveManagedTeacher(c.Request.Context(), managerID, teacherID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// UpdateMyManagedCategories updates the categories managed by the current manager
// @Summary Update my managed categories
// @Description Updates the list of categories managed by the current manager
// @Tags managers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body UpdateManagedResourcesRequest true "Update Categories Request"
// @Success 204 "No Content"
// @Router /api/v1/managers/me/categories [put]
func (h *ManagerHandler) UpdateMyManagedCategories(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get the manager by user ID first
	m, err := h.service.GetByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	var req UpdateManagedResourcesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove all existing categories first
	if m.ManagesCategories != nil {
		for _, catID := range m.ManagesCategories {
			h.service.RemoveManagedCategory(c.Request.Context(), m.ID, catID)
		}
	}

	// Add new categories
	if req.Categories != nil {
		for _, catID := range req.Categories {
			if err := h.service.AddManagedCategory(c.Request.Context(), m.ID, catID); err != nil {
				status := getStatusCode(err)
				c.JSON(status, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusNoContent, nil)
}

// UpdateMyManagedTeachers updates the teachers managed by the current manager
// @Summary Update my managed teachers
// @Description Updates the list of teachers managed by the current manager
// @Tags managers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body UpdateManagedResourcesRequest true "Update Teachers Request"
// @Success 204 "No Content"
// @Router /api/v1/managers/me/teachers [put]
func (h *ManagerHandler) UpdateMyManagedTeachers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get the manager by user ID first
	m, err := h.service.GetByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	var req UpdateManagedResourcesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove all existing teachers first
	if m.ManagesTeachers != nil {
		for _, teacherID := range m.ManagesTeachers {
			h.service.RemoveManagedTeacher(c.Request.Context(), m.ID, teacherID)
		}
	}

	// Add new teachers
	if req.Teachers != nil {
		for _, teacherID := range req.Teachers {
			if err := h.service.AddManagedTeacher(c.Request.Context(), m.ID, teacherID); err != nil {
				status := getStatusCode(err)
				c.JSON(status, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusNoContent, nil)
}
