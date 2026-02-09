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

func (h *CourseCategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
