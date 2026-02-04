package notification

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/notification"
	"github.com/google/uuid"
)

type Service struct {
	repo notification.Repository
}

func NewService(repo notification.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input *notification.CreateNotificationInput) (*notification.Notification, error) {
	n := &notification.Notification{
		ID:        uuid.New().String(),
		UserID:    input.UserID,
		Title:     input.Title,
		Message:   input.Message,
		Read:      false,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) GetByID(id string) (*notification.Notification, error) {
	return s.repo.GetByID(id)
}

func (s *Service) ListByUser(userID string, skip, take int) ([]*notification.Notification, error) {
	return s.repo.ListByUser(userID, skip, take)
}

func (s *Service) MarkAsRead(id string) error {
	return s.repo.MarkAsRead(id)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
