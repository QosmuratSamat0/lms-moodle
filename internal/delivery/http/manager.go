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

func NewManagerHandler(service *managerUC.Service) *ManagerHandler {
	return &ManagerHandler{service: service}
}

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

func (h *ManagerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteManager(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

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

func (h *ManagerHandler) AddManagedCategory(c *gin.Context) {
	managerID := c.Param("id")
	var req struct {
		CategoryID string `json:"category_id" binding:"required"`
	}
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

func (h *ManagerHandler) RemoveManagedCategory(c *gin.Context) {
	managerID := c.Param("id")
	categoryID := c.Param("categoryID")

	if err := h.service.RemoveManagedCategory(c.Request.Context(), managerID, categoryID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *ManagerHandler) AddManagedTeacher(c *gin.Context) {
	managerID := c.Param("id")
	var req struct {
		TeacherID string `json:"teacher_id" binding:"required"`
	}
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

func (h *ManagerHandler) RemoveManagedTeacher(c *gin.Context) {
	managerID := c.Param("id")
	teacherID := c.Param("teacherID")

	if err := h.service.RemoveManagedTeacher(c.Request.Context(), managerID, teacherID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
