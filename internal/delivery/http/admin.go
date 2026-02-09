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

func (h *AdminHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteAdmin(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
