package chat

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/chat"
	"github.com/google/uuid"
)

type Service struct {
	repo chat.Repository
}

func NewService(repo chat.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendMessage(input *chat.CreateMessageInput) (*chat.Message, error) {
	m := &chat.Message{
		ID:        uuid.New().String(),
		SenderID:  input.SenderID,
		CourseID:  input.CourseID,
		Content:   input.Content,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) GetByID(id string) (*chat.Message, error) {
	return s.repo.GetByID(id)
}

func (s *Service) ListByCourse(courseID string, skip, take int) ([]*chat.Message, error) {
	return s.repo.ListByCourse(courseID, skip, take)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
