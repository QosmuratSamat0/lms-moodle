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

// Upload records a file upload
// @Summary Record upload
// @Description Records file upload metadata
// @Tags uploads
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body struct{UserID string `json:"user_id" binding:"required"`; FileName string `json:"file_name" binding:"required"`; FileURL string `json:"file_url" binding:"required"`; FileSize int64 `json:"file_size" binding:"required"`; MimeType string `json:"mime_type" binding:"required"`} true "Upload Request"
// @Success 201 {object} upload.Upload "Recorded upload"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/uploads [post]
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

// GetByID returns an upload record by ID
// @Summary Get upload by ID
// @Description Returns upload metadata by ID
// @Tags uploads
// @Security BearerAuth
// @Produce json
// @Param id path string true "Upload ID"
// @Success 200 {object} upload.Upload "Upload details"
// @Router /api/v1/uploads/{id} [get]
func (h *UploadHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	u, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "upload not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// ListByUser returns upload records for a specific user
// @Summary List user uploads
// @Description Returns a paginated list of uploads for a user
// @Tags uploads
// @Security BearerAuth
// @Produce json
// @Param userID path string true "User ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} upload.Upload "Uploads list"
// @Router /api/v1/uploads/user/{userID} [get]
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

// Delete deletes an upload record
// @Summary Delete upload
// @Description Deletes an upload record by ID (Admin only)
// @Tags uploads
// @Security BearerAuth
// @Produce json
// @Param id path string true "Upload ID"
// @Success 204 "No content"
// @Router /api/v1/uploads/{id} [delete]
func (h *UploadHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
