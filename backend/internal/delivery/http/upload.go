package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/upload"
	uploadUC "github.com/ap1-final-mini-moodle/internal/usecase/upload"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *uploadUC.Service
}

func NewUploadHandler(service *uploadUC.Service) *UploadHandler {
	return &UploadHandler{service: service}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		FileName string `json:"file_name" binding:"required"`
		FileURL  string `json:"file_url" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
		MimeType string `json:"mime_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.service.Upload(&upload.CreateUploadInput{
		UserID:   req.UserID,
		FileName: req.FileName,
		FileURL:  req.FileURL,
		FileSize: req.FileSize,
		MimeType: req.MimeType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *UploadHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	u, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "upload not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UploadHandler) ListByUser(c *gin.Context) {
	userID := c.Param("userID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uploads, err := h.service.ListByUser(userID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, uploads)
}

func (h *UploadHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
