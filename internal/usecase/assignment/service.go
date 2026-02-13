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
	// Ensure max_points is at least 1 to satisfy database constraint
	maxPoints := input.MaxPoints
	if maxPoints <= 0 {
		maxPoints = 1
	}

	a := &assignment.Assignment{
		ID:                 uuid.New().String(),
		CourseID:           input.CourseID,
		Title:              input.Title,
		Description:        input.Description,
		MaxPoints:          maxPoints,
		DueAt:              input.DueAt,
		AllowLate:          input.AllowLate,
		CreatedByTeacherID: input.CreatedByTeacherID,
		GradingCategory:    input.GradingCategory,
		WeightPercentage:   input.WeightPercentage,
		FileURL:            input.FileURL,
		CreatedAt:          time.Now(),
	}
	if a.GradingCategory == "" {
		a.GradingCategory = "register_midterm"
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
	if input.DueAt != nil {
		a.DueAt = input.DueAt
	}
	if input.AllowLate != nil {
		a.AllowLate = *input.AllowLate
	}
	if input.GradingCategory != nil {
		a.GradingCategory = *input.GradingCategory
	}
	if input.WeightPercentage != nil {
		a.WeightPercentage = *input.WeightPercentage
	}
	if input.FileURL != nil {
		a.FileURL = input.FileURL
	}
	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
