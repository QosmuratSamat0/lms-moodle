package enrollment

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/enrollment"
	"github.com/google/uuid"
)

type Service struct {
	repo enrollment.Repository
}

func NewService(repo enrollment.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Enroll(input *enrollment.CreateEnrollmentInput) (*enrollment.Enrollment, error) {
	e := &enrollment.Enrollment{
		ID:         uuid.New().String(),
		CourseID:   input.CourseID,
		StudentID:  input.StudentID,
		EnrolledAt: time.Now(),
		Status:     "active",
	}
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) GetByID(id string) (*enrollment.Enrollment, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetByCourseAndStudent(courseID, studentID string) (*enrollment.Enrollment, error) {
	return s.repo.GetByCourseAndStudent(courseID, studentID)
}

func (s *Service) ListByCourse(courseID string, skip, take int) ([]*enrollment.Enrollment, error) {
	return s.repo.ListByCourse(courseID, skip, take)
}

func (s *Service) ListByStudent(studentID string, skip, take int) ([]*enrollment.Enrollment, error) {
	return s.repo.ListByStudent(studentID, skip, take)
}

func (s *Service) Remove(id string) error {
	return s.repo.Delete(id)
}
