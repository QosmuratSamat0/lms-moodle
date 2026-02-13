package attendance

import (
	"fmt"
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

func (s *Service) CreateSession(input *attendance.CreateSessionInput) (*attendance.AttendanceSession, error) {
	session := &attendance.AttendanceSession{
		ID:                 uuid.New().String(),
		CourseID:           input.CourseID,
		StartsAt:           input.StartsAt,
		EndsAt:             input.EndsAt,
		CreatedByTeacherID: &input.TeacherID,
		CreatedAt:          time.Now(),
	}
	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) GetSession(id string) (*attendance.AttendanceSession, error) {
	return s.repo.GetSession(id)
}

func (s *Service) ListSessionsByCourse(courseID string) ([]*attendance.AttendanceSession, error) {
	return s.repo.ListSessionsByCourse(courseID)
}

func (s *Service) DeleteSession(id string) error {
	return s.repo.DeleteSession(id)
}

func (s *Service) MarkAttendance(input *attendance.MarkAttendanceInput) (*attendance.AttendanceMark, error) {
	validStatuses := map[string]bool{"present": true, "absent": true, "late": true, "excused": true}
	if !validStatuses[input.Status] {
		return nil, fmt.Errorf("invalid status: %s", input.Status)
	}
	mark := &attendance.AttendanceMark{
		ID:        uuid.New().String(),
		SessionID: input.SessionID,
		StudentID: input.StudentID,
		Status:    input.Status,
		MarkedAt:  time.Now(),
	}
	if err := s.repo.UpsertMark(mark); err != nil {
		return nil, err
	}
	return mark, nil
}

func (s *Service) BulkMark(input *attendance.BulkMarkInput) error {
	validStatuses := map[string]bool{"present": true, "absent": true, "late": true, "excused": true}
	var marks []*attendance.AttendanceMark
	for _, m := range input.Marks {
		if !validStatuses[m.Status] {
			return fmt.Errorf("invalid status for student %s: %s", m.StudentID, m.Status)
		}
		marks = append(marks, &attendance.AttendanceMark{
			ID:        uuid.New().String(),
			SessionID: input.SessionID,
			StudentID: m.StudentID,
			Status:    m.Status,
			MarkedAt:  time.Now(),
		})
	}
	return s.repo.BulkUpsertMarks(marks)
}

func (s *Service) GetMarksBySession(sessionID string) ([]*attendance.AttendanceMark, error) {
	return s.repo.GetMarksBySession(sessionID)
}

func (s *Service) GetStudentAttendance(studentID, courseID string) ([]*attendance.AttendanceMark, error) {
	return s.repo.GetStudentAttendance(studentID, courseID)
}

func (s *Service) GetStudentSummary(studentID, courseID string) (*attendance.StudentAttendanceSummary, error) {
	return s.repo.GetStudentSummary(studentID, courseID)
}
