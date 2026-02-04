package upload

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/upload"
	"github.com/google/uuid"
)

type Service struct {
	repo upload.Repository
}

func NewService(repo upload.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upload(input *upload.CreateUploadInput) (*upload.Upload, error) {
	u := &upload.Upload{
		ID:        uuid.New().String(),
		UserID:    input.UserID,
		FileName:  input.FileName,
		FileURL:   input.FileURL,
		FileSize:  input.FileSize,
		MimeType:  input.MimeType,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) GetByID(id string) (*upload.Upload, error) {
	return s.repo.GetByID(id)
}

func (s *Service) ListByUser(userID string, skip, take int) ([]*upload.Upload, error) {
	return s.repo.ListByUser(userID, skip, take)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
