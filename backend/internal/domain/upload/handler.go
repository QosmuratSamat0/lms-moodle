package upload

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles upload-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new upload handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Upload handles file upload
// @Summary Upload a file
// @Tags uploads
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param type formData string true "Attachment type" Enums(submission, assignment, course, profile, chat)
// @Param reference_id formData string false "Reference ID (submission_id, assignment_id, etc.)"
// @Success 201 {object} AttachmentResponse
// @Failure 400 {object} utils.Response
// @Router /uploads [post]
func (h *Handler) Upload(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	// Get attachment type
	attachmentType := AttachmentType(c.PostForm("type"))
	if attachmentType == "" {
		utils.BadRequest(c, "type is required")
		return
	}

	// Validate attachment type
	validTypes := map[AttachmentType]bool{
		AttachmentTypeSubmission: true,
		AttachmentTypeAssignment: true,
		AttachmentTypeCourse:     true,
		AttachmentTypeProfile:    true,
		AttachmentTypeChat:       true,
	}
	if !validTypes[attachmentType] {
		utils.BadRequest(c, "invalid attachment type")
		return
	}

	// Get optional reference ID
	var referenceID *uuid.UUID
	if refIDStr := c.PostForm("reference_id"); refIDStr != "" {
		refID, err := uuid.Parse(refIDStr)
		if err != nil {
			utils.BadRequest(c, "invalid reference_id")
			return
		}
		referenceID = &refID
	}

	attachment, err := h.service.Upload(c.Request.Context(), userID, file, header, attachmentType, referenceID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, attachment)
}

// UploadFromURL handles file upload from URL
// @Summary Upload file from URL
// @Tags uploads
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body UploadFromURLRequest true "Upload details"
// @Success 201 {object} AttachmentResponse
// @Failure 400 {object} utils.Response
// @Router /uploads/url [post]
func (h *Handler) UploadFromURL(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req UploadFromURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	attachment, err := h.service.UploadFromURL(c.Request.Context(), userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, attachment)
}

// GetByID retrieves an attachment by ID
// @Summary Get attachment by ID
// @Tags uploads
// @Security BearerAuth
// @Produce json
// @Param id path string true "Attachment ID"
// @Success 200 {object} AttachmentResponse
// @Failure 404 {object} utils.Response
// @Router /uploads/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid attachment ID")
		return
	}

	attachment, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, attachment)
}

// Delete deletes an attachment
// @Summary Delete attachment
// @Tags uploads
// @Security BearerAuth
// @Param id path string true "Attachment ID"
// @Success 204
// @Failure 403,404 {object} utils.Response
// @Router /uploads/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid attachment ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	if err := h.service.Delete(c.Request.Context(), id, userID, role); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListMyUploads lists the current user's uploads
// @Summary List my uploads
// @Tags uploads
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} AttachmentListResponse
// @Router /uploads/my [get]
func (h *Handler) ListMyUploads(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListMyUploads(c.Request.Context(), userID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// ListByReference lists attachments by reference
// @Summary List attachments by reference
// @Tags uploads
// @Security BearerAuth
// @Produce json
// @Param type path string true "Attachment type" Enums(submission, assignment, course, profile, chat)
// @Param reference_id path string true "Reference ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} AttachmentListResponse
// @Router /uploads/reference/{type}/{reference_id} [get]
func (h *Handler) ListByReference(c *gin.Context) {
	refType := AttachmentType(c.Param("type"))
	refIDStr := c.Param("reference_id")

	refID, err := uuid.Parse(refIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid reference_id")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListByReference(c.Request.Context(), refType, refID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// GetAllowedTypes returns allowed file types
// @Summary Get allowed file types
// @Tags uploads
// @Produce json
// @Success 200 {object} AllowedTypesResponse
// @Router /uploads/allowed-types [get]
func (h *Handler) GetAllowedTypes(c *gin.Context) {
	utils.OK(c, h.service.GetAllowedTypes())
}
