package grade

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/grade"
	"github.com/google/uuid"
)

type Service struct {
	repo grade.Repository
}

func NewService(repo grade.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Grade(input *grade.CreateGradeInput) (*grade.Grade, error) {
	g := &grade.Grade{
		ID:           uuid.New().String(),
		SubmissionID: input.SubmissionID,
		Score:        input.Score,
		Feedback:     input.Feedback,
		GradedBy:     input.GradedBy,
		GradedAt:     time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.repo.Create(g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *Service) GetByID(id string) (*grade.Grade, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetBySubmission(submissionID string) (*grade.Grade, error) {
	return s.repo.GetBySubmission(submissionID)
}

func (s *Service) Update(id string, score int, feedback string) (*grade.Grade, error) {
	g, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	g.Score = score
	g.Feedback = feedback
	g.UpdatedAt = time.Now()
	if err := s.repo.Update(g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
