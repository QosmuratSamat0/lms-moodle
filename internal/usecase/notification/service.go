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
	notifType := input.Type
	if notifType == "" {
		notifType = "info"
	}
	n := &notification.Notification{
		ID:        uuid.New().String(),
		UserID:    input.UserID,
		Type:      notifType,
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

// NotifyMany sends the same notification to multiple users
func (s *Service) NotifyMany(userIDs []string, notifType, title, message string) {
	var notifications []*notification.Notification
	now := time.Now()
	for _, uid := range userIDs {
		notifications = append(notifications, &notification.Notification{
			ID:        uuid.New().String(),
			UserID:    uid,
			Type:      notifType,
			Title:     title,
			Message:   message,
			Read:      false,
			CreatedAt: now,
		})
	}
	// Best-effort: don't fail the main action if notifications fail
	_ = s.repo.CreateBulk(notifications)
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

func (s *Service) MarkAllAsRead(userID string) error {
	return s.repo.MarkAllAsRead(userID)
}

func (s *Service) UnreadCount(userID string) (int, error) {
	return s.repo.UnreadCount(userID)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
