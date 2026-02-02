package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// AllowedFileTypes defines permitted file extensions and their MIME types
var AllowedFileTypes = map[string][]string{
	// Images
	".jpg":  {"image/jpeg"},
	".jpeg": {"image/jpeg"},
	".png":  {"image/png"},
	".gif":  {"image/gif"},
	".webp": {"image/webp"},
	// Documents
	".pdf":  {"application/pdf"},
	".doc":  {"application/msword"},
	".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	".xls":  {"application/vnd.ms-excel"},
	".xlsx": {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
	".ppt":  {"application/vnd.ms-powerpoint"},
	".pptx": {"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
	".txt":  {"text/plain"},
	// Archives
	".zip": {"application/zip"},
}

// Config holds Cloudinary configuration
type Config struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadPreset string
	Folder       string
	MaxFileSize  int64 // in bytes
}

// CloudinaryUploader handles file uploads to Cloudinary
type CloudinaryUploader struct {
	cld         *cloudinary.Cloudinary
	config      Config
	maxFileSize int64
}

// UploadResult represents the result of a file upload
type UploadResult struct {
	PublicID     string    `json:"public_id"`
	URL          string    `json:"url"`
	SecureURL    string    `json:"secure_url"`
	OriginalName string    `json:"original_name"`
	Format       string    `json:"format"`
	ResourceType string    `json:"resource_type"`
	Size         int64     `json:"size"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewCloudinaryUploader creates a new Cloudinary uploader instance
func NewCloudinaryUploader(cfg Config) (*CloudinaryUploader, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudinary client: %w", err)
	}

	maxSize := cfg.MaxFileSize
	if maxSize == 0 {
		maxSize = 10 * 1024 * 1024 // 10MB default
	}

	return &CloudinaryUploader{
		cld:         cld,
		config:      cfg,
		maxFileSize: maxSize,
	}, nil
}

// ValidateFile checks if the file is allowed based on extension and size
func (u *CloudinaryUploader) ValidateFile(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > u.maxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", header.Size, u.maxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if _, ok := AllowedFileTypes[ext]; !ok {
		return fmt.Errorf("file type %s is not allowed", ext)
	}

	return nil
}

// Upload uploads a file to Cloudinary
func (u *CloudinaryUploader) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, subfolder string) (*UploadResult, error) {
	// Validate file first
	if err := u.ValidateFile(header); err != nil {
		return nil, err
	}

	// Determine resource type based on extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	resourceType := "auto"
	if strings.HasPrefix(ext, ".jpg") || strings.HasPrefix(ext, ".jpeg") ||
		strings.HasPrefix(ext, ".png") || strings.HasPrefix(ext, ".gif") ||
		strings.HasPrefix(ext, ".webp") {
		resourceType = "image"
	} else {
		resourceType = "raw"
	}

	// Build folder path
	folder := u.config.Folder
	if subfolder != "" {
		folder = fmt.Sprintf("%s/%s", folder, subfolder)
	}

	// Generate unique public ID
	timestamp := time.Now().Unix()
	baseName := strings.TrimSuffix(header.Filename, ext)
	publicID := fmt.Sprintf("%s/%d_%s", folder, timestamp, sanitizeFilename(baseName))

	// Upload to Cloudinary
	uploadParams := uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: resourceType,
		Folder:       folder,
	}

	if u.config.UploadPreset != "" {
		uploadParams.UploadPreset = u.config.UploadPreset
	}

	result, err := u.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &UploadResult{
		PublicID:     result.PublicID,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		OriginalName: header.Filename,
		Format:       result.Format,
		ResourceType: result.ResourceType,
		Size:         int64(result.Bytes),
		Width:        result.Width,
		Height:       result.Height,
		CreatedAt:    time.Now(),
	}, nil
}

// UploadFromURL uploads a file from URL to Cloudinary
func (u *CloudinaryUploader) UploadFromURL(ctx context.Context, url string, subfolder string) (*UploadResult, error) {
	folder := u.config.Folder
	if subfolder != "" {
		folder = fmt.Sprintf("%s/%s", folder, subfolder)
	}

	uploadParams := uploader.UploadParams{
		Folder:       folder,
		ResourceType: "auto",
	}

	result, err := u.cld.Upload.Upload(ctx, url, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload from URL: %w", err)
	}

	return &UploadResult{
		PublicID:     result.PublicID,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		OriginalName: filepath.Base(url),
		Format:       result.Format,
		ResourceType: result.ResourceType,
		Size:         int64(result.Bytes),
		Width:        result.Width,
		Height:       result.Height,
		CreatedAt:    time.Now(),
	}, nil
}

// Delete removes a file from Cloudinary
func (u *CloudinaryUploader) Delete(ctx context.Context, publicID string, resourceType string) error {
	if resourceType == "" {
		resourceType = "image"
	}

	_, err := u.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetAllowedExtensions returns a list of allowed file extensions
func GetAllowedExtensions() []string {
	exts := make([]string, 0, len(AllowedFileTypes))
	for ext := range AllowedFileTypes {
		exts = append(exts, ext)
	}
	return exts
}

// sanitizeFilename removes special characters from filename
func sanitizeFilename(name string) string {
	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")
	// Remove special characters
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	var result strings.Builder
	for _, char := range name {
		if strings.ContainsRune(allowed, char) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// MaxFileSize returns the maximum allowed file size in bytes
func (u *CloudinaryUploader) MaxFileSize() int64 {
	return u.maxFileSize
}
