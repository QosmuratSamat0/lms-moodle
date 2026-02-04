package assignment

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/assignment"
	"github.com/google/uuid"
)

type Service struct {
	repo assignment.Repository
}

func NewService(repo assignment.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input *assignment.CreateAssignmentInput) (*assignment.Assignment, error) {
	a := &assignment.Assignment{
		ID:          uuid.New().String(),
		CourseID:    input.CourseID,
		Title:       input.Title,
		Description: input.Description,
		MaxPoints:   input.MaxPoints,
		DueDate:     input.DueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetByID(id string) (*assignment.Assignment, error) {
	return s.repo.GetByID(id)
}

func (s *Service) ListByCourse(courseID string, skip, take int) ([]*assignment.Assignment, error) {
	return s.repo.ListByCourse(courseID, skip, take)
}

func (s *Service) Update(id string, input *assignment.UpdateAssignmentInput) (*assignment.Assignment, error) {
	a, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Title != nil {
		a.Title = *input.Title
	}
	if input.Description != nil {
		a.Description = *input.Description
	}
	if input.MaxPoints != nil {
		a.MaxPoints = *input.MaxPoints
	}
	if input.DueDate != nil {
		a.DueDate = *input.DueDate
	}
	a.UpdatedAt = time.Now()
	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
