package upload

import "github.com/google/uuid"

// UploadRequest represents a request to upload a file
// Note: This is not used directly in JSON binding since we use multipart/form-data
type UploadRequest struct {
	Type        AttachmentType `form:"type" validate:"required,oneof=submission assignment course profile chat"`
	ReferenceID *uuid.UUID     `form:"reference_id" validate:"omitempty"`
}

// UploadFromURLRequest represents a request to upload from URL
type UploadFromURLRequest struct {
	URL         string         `json:"url" validate:"required,url"`
	Type        AttachmentType `json:"type" validate:"required,oneof=submission assignment course profile chat"`
	ReferenceID *uuid.UUID     `json:"reference_id,omitempty"`
}

// AttachmentResponse represents an attachment in API responses
type AttachmentResponse struct {
	ID                string  `json:"id"`
	UploadedBy        string  `json:"uploaded_by"`
	Type              string  `json:"type"`
	ReferenceID       *string `json:"reference_id,omitempty"`
	PublicID          string  `json:"public_id"`
	URL               string  `json:"url"`
	SecureURL         string  `json:"secure_url"`
	OriginalName      string  `json:"original_name"`
	Format            string  `json:"format"`
	ResourceType      string  `json:"resource_type"`
	Size              int64   `json:"size"`
	Width             *int    `json:"width,omitempty"`
	Height            *int    `json:"height,omitempty"`
	UploaderFirstName *string `json:"uploader_first_name,omitempty"`
	UploaderLastName  *string `json:"uploader_last_name,omitempty"`
	UploaderEmail     string  `json:"uploader_email,omitempty"`
	CreatedAt         string  `json:"created_at"`
}

// AttachmentListResponse represents a paginated list of attachments
type AttachmentListResponse struct {
	Attachments []AttachmentResponse `json:"attachments"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
	TotalPages  int                  `json:"total_pages"`
}

// UploadSignatureResponse represents a signed upload URL response
type UploadSignatureResponse struct {
	Signature  string `json:"signature"`
	Timestamp  int64  `json:"timestamp"`
	CloudName  string `json:"cloud_name"`
	APIKey     string `json:"api_key"`
	Folder     string `json:"folder"`
	UploadURL  string `json:"upload_url"`
}

// AllowedTypesResponse returns allowed file types
type AllowedTypesResponse struct {
	Extensions  []string `json:"extensions"`
	MaxFileSize int64    `json:"max_file_size"`
}
