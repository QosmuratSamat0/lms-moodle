package upload

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	cloudupload "github.com/ap1-final-mini-moodle/pkg/upload"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the upload service interface
type Service interface {
	Upload(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader, attachmentType AttachmentType, referenceID *uuid.UUID) (*AttachmentResponse, error)
	UploadFromURL(ctx context.Context, userID uuid.UUID, req *UploadFromURLRequest) (*AttachmentResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*AttachmentResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string) error
	ListByReference(ctx context.Context, refType AttachmentType, refID uuid.UUID, page, limit int) (*AttachmentListResponse, error)
	ListMyUploads(ctx context.Context, userID uuid.UUID, page, limit int) (*AttachmentListResponse, error)
	GetAllowedTypes() *AllowedTypesResponse
}

type service struct {
	repo     Repository
	uploader *cloudupload.CloudinaryUploader
}

// NewService creates a new upload service
func NewService(repo Repository, uploader *cloudupload.CloudinaryUploader) Service {
	return &service{
		repo:     repo,
		uploader: uploader,
	}
}

func (s *service) Upload(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader, attachmentType AttachmentType, referenceID *uuid.UUID) (*AttachmentResponse, error) {
	// Upload to Cloudinary
	subfolder := string(attachmentType)
	result, err := s.uploader.Upload(ctx, file, header, subfolder)
	if err != nil {
		return nil, errorx.Wrap(err, "upload file to cloudinary")
	}

	// Save attachment metadata to database
	attachment := &Attachment{
		UploadedBy:   userID,
		Type:         attachmentType,
		ReferenceID:  referenceID,
		PublicID:     result.PublicID,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		OriginalName: result.OriginalName,
		Format:       result.Format,
		ResourceType: result.ResourceType,
		Size:         result.Size,
	}

	if result.Width > 0 {
		attachment.Width = &result.Width
	}
	if result.Height > 0 {
		attachment.Height = &result.Height
	}

	if err := s.repo.Create(ctx, attachment); err != nil {
		// Try to clean up the uploaded file from Cloudinary
		_ = s.uploader.Delete(ctx, result.PublicID, result.ResourceType)
		return nil, errorx.Wrap(err, "save attachment metadata")
	}

	// Fetch with user details
	detailed, err := s.repo.GetByID(ctx, attachment.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get attachment details")
	}

	return s.toResponse(detailed), nil
}

func (s *service) UploadFromURL(ctx context.Context, userID uuid.UUID, req *UploadFromURLRequest) (*AttachmentResponse, error) {
	// Upload to Cloudinary from URL
	subfolder := string(req.Type)
	result, err := s.uploader.UploadFromURL(ctx, req.URL, subfolder)
	if err != nil {
		return nil, errorx.Wrap(err, "upload from URL to cloudinary")
	}

	// Save attachment metadata
	attachment := &Attachment{
		UploadedBy:   userID,
		Type:         req.Type,
		ReferenceID:  req.ReferenceID,
		PublicID:     result.PublicID,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		OriginalName: result.OriginalName,
		Format:       result.Format,
		ResourceType: result.ResourceType,
		Size:         result.Size,
	}

	if result.Width > 0 {
		attachment.Width = &result.Width
	}
	if result.Height > 0 {
		attachment.Height = &result.Height
	}

	if err := s.repo.Create(ctx, attachment); err != nil {
		_ = s.uploader.Delete(ctx, result.PublicID, result.ResourceType)
		return nil, errorx.Wrap(err, "save attachment metadata")
	}

	detailed, err := s.repo.GetByID(ctx, attachment.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get attachment details")
	}

	return s.toResponse(detailed), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*AttachmentResponse, error) {
	attachment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("attachment")
		}
		return nil, errorx.Wrap(err, "get attachment")
	}
	return s.toResponse(attachment), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string) error {
	attachment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("attachment")
		}
		return errorx.Wrap(err, "get attachment")
	}

	// Check ownership (admin can delete any)
	if role != "admin" && attachment.UploadedBy != userID {
		return errorx.NewForbiddenError("you don't have permission to delete this attachment")
	}

	// Delete from Cloudinary
	if err := s.uploader.Delete(ctx, attachment.PublicID, attachment.ResourceType); err != nil {
		// Log but don't fail - still delete from DB
	}

	// Delete from database
	if err := s.repo.Delete(ctx, id); err != nil {
		return errorx.Wrap(err, "delete attachment")
	}

	return nil
}

func (s *service) ListByReference(ctx context.Context, refType AttachmentType, refID uuid.UUID, page, limit int) (*AttachmentListResponse, error) {
	offset := (page - 1) * limit
	attachments, total, err := s.repo.ListByReference(ctx, refType, refID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list attachments")
	}
	return s.toListResponse(attachments, total, page, limit), nil
}

func (s *service) ListMyUploads(ctx context.Context, userID uuid.UUID, page, limit int) (*AttachmentListResponse, error) {
	offset := (page - 1) * limit
	attachments, total, err := s.repo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list my uploads")
	}
	return s.toListResponse(attachments, total, page, limit), nil
}

func (s *service) GetAllowedTypes() *AllowedTypesResponse {
	return &AllowedTypesResponse{
		Extensions:  cloudupload.GetAllowedExtensions(),
		MaxFileSize: s.uploader.MaxFileSize(),
	}
}

func (s *service) toResponse(att *AttachmentWithUser) *AttachmentResponse {
	resp := &AttachmentResponse{
		ID:                att.ID.String(),
		UploadedBy:        att.UploadedBy.String(),
		Type:              string(att.Type),
		PublicID:          att.PublicID,
		URL:               att.URL,
		SecureURL:         att.SecureURL,
		OriginalName:      att.OriginalName,
		Format:            att.Format,
		ResourceType:      att.ResourceType,
		Size:              att.Size,
		Width:             att.Width,
		Height:            att.Height,
		UploaderFirstName: att.UploaderFirstName,
		UploaderLastName:  att.UploaderLastName,
		UploaderEmail:     att.UploaderEmail,
		CreatedAt:         att.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if att.ReferenceID != nil {
		refStr := att.ReferenceID.String()
		resp.ReferenceID = &refStr
	}

	return resp
}

func (s *service) toListResponse(attachments []AttachmentWithUser, total int64, page, limit int) *AttachmentListResponse {
	responses := make([]AttachmentResponse, 0, len(attachments))
	for _, att := range attachments {
		responses = append(responses, *s.toResponse(&att))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &AttachmentListResponse{
		Attachments: responses,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
	}
}
