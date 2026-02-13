package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ap1-final-mini-moodle/internal/domain/upload"
	uploadUC "github.com/ap1-final-mini-moodle/internal/usecase/upload"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	service   *uploadUC.Service
	uploadDir string
}

func NewUploadHandler(service *uploadUC.Service) *UploadHandler {
	dir := "./uploads"
	os.MkdirAll(dir, 0755)
	return &UploadHandler{service: service, uploadDir: dir}
}

// Upload handles actual file upload via multipart form
// @Summary Upload a file
// @Description Uploads a file and returns its metadata with a download URL
// @Tags uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} map[string]interface{} "Upload result"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/uploads [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided: " + err.Error()})
		return
	}
	defer file.Close()

	// Get user ID from JWT
	userID := ""
	if uid, exists := c.Get("userID"); exists {
		userID = uid.(string)
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	safeExt := strings.ToLower(ext)
	uniqueName := uuid.New().String() + safeExt

	// Save file to disk
	destPath := filepath.Join(h.uploadDir, uniqueName)
	dest, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	// Build the download URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	fileURL := fmt.Sprintf("%s://%s/api/v1/uploads/files/%s", scheme, host, uniqueName)

	// Detect MIME type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Record in database
	u, err := h.service.Upload(&upload.CreateUploadInput{
		UserID:   userID,
		FileName: header.Filename,
		FileURL:  fileURL,
		FileSize: header.Size,
		MimeType: mimeType,
	})
	if err != nil {
		// File saved but DB record failed — still return URL
		c.JSON(http.StatusCreated, gin.H{
			"id":            "",
			"url":           fileURL,
			"secure_url":    fileURL,
			"original_name": header.Filename,
			"size":          header.Size,
			"format":        safeExt,
			"resource_type": mimeType,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            u.ID,
		"url":           fileURL,
		"secure_url":    fileURL,
		"original_name": header.Filename,
		"size":          header.Size,
		"format":        safeExt,
		"resource_type": mimeType,
		"file_url":      fileURL,
		"file_name":     header.Filename,
		"created_at":    u.CreatedAt,
	})
}

// ServeFile serves an uploaded file by filename
func (h *UploadHandler) ServeFile(c *gin.Context) {
	filename := c.Param("filename")
	// Prevent directory traversal
	filename = filepath.Base(filename)
	filePath := filepath.Join(h.uploadDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
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
