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

func (h *CategoryManagerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

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
