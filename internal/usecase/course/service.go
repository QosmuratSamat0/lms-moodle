package course

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/course"
	"github.com/google/uuid"
)

type Service struct {
	repo course.Repository
}

func NewService(repo course.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input *course.CreateCourseInput) (*course.Course, error) {
	c := &course.Course{
		ID:          uuid.New().String(),
		Code:        input.Code,
		Title:       input.Title,
		Description: input.Description,
		TeacherID:   input.TeacherID,
		MaxPoints:   input.MaxPoints,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetByID(id string) (*course.Course, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetByCode(code string) (*course.Course, error) {
	return s.repo.GetByCode(code)
}

func (s *Service) Update(id string, input *course.UpdateCourseInput) (*course.Course, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Title != nil {
		c.Title = *input.Title
	}
	if input.Description != nil {
		c.Description = *input.Description
	}
	if input.MaxPoints != nil {
		c.MaxPoints = *input.MaxPoints
	}
	if input.Active != nil {
		c.Active = *input.Active
	}
	c.UpdatedAt = time.Now()
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) List(skip, take int) ([]*course.Course, error) {
	return s.repo.List(skip, take)
}

func (s *Service) ListByTeacher(teacherID string, skip, take int) ([]*course.Course, error) {
	return s.repo.ListByTeacher(teacherID, skip, take)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
