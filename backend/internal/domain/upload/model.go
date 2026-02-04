package upload

import (
	"time"

	"github.com/google/uuid"
)

// AttachmentType represents the type/category of attachment
type AttachmentType string

const (
	AttachmentTypeSubmission AttachmentType = "submission"
	AttachmentTypeAssignment AttachmentType = "assignment"
	AttachmentTypeCourse     AttachmentType = "course"
	AttachmentTypeProfile    AttachmentType = "profile"
	AttachmentTypeChat       AttachmentType = "chat"
)

// Attachment represents a file attachment in the system
type Attachment struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	UploadedBy   uuid.UUID      `json:"uploaded_by" db:"uploaded_by"`
	Type         AttachmentType `json:"type" db:"type"`
	ReferenceID  *uuid.UUID     `json:"reference_id" db:"reference_id"` // ID of related entity (submission, assignment, etc.)
	PublicID     string         `json:"public_id" db:"public_id"`       // Cloudinary public ID
	URL          string         `json:"url" db:"url"`
	SecureURL    string         `json:"secure_url" db:"secure_url"`
	OriginalName string         `json:"original_name" db:"original_name"`
	Format       string         `json:"format" db:"format"`
	ResourceType string         `json:"resource_type" db:"resource_type"`
	Size         int64          `json:"size" db:"size"`
	Width        *int           `json:"width,omitempty" db:"width"`
	Height       *int           `json:"height,omitempty" db:"height"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}

// AttachmentWithUser includes uploader information
type AttachmentWithUser struct {
	Attachment
	UploaderFirstName *string `json:"uploader_first_name,omitempty" db:"uploader_first_name"`
	UploaderLastName  *string `json:"uploader_last_name,omitempty" db:"uploader_last_name"`
	UploaderEmail     string  `json:"uploader_email" db:"uploader_email"`
}
