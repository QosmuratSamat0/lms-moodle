package attendance

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/attendance"
	"github.com/google/uuid"
)

type Service struct {
	repo attendance.Repository
}

func NewService(repo attendance.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Record(input *attendance.CreateAttendanceInput) (*attendance.Attendance, error) {
	a := &attendance.Attendance{
		ID:        uuid.New().String(),
		CourseID:  input.CourseID,
		StudentID: input.StudentID,
		Date:      input.Date,
		Present:   input.Present,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) ListByCourse(courseID string, skip, take int) ([]*attendance.Attendance, error) {
	return s.repo.ListByCourse(courseID, skip, take)
}

func (s *Service) ListByStudent(studentID string, skip, take int) ([]*attendance.Attendance, error) {
	return s.repo.ListByStudent(studentID, skip, take)
}

func (s *Service) GetByStudentAndDate(studentID string, date time.Time) (*attendance.Attendance, error) {
	return s.repo.GetByStudentAndDate(studentID, date)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
