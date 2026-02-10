package student

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/student"
	"github.com/ap1-final-mini-moodle/internal/domain/user"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo     student.Repository
	userRepo user.Repository
}

func NewService(repo student.Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) CreateStudent(ctx context.Context, input *student.CreateStudentInput) (*student.Student, error) {
	if input.UserID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.StudentCode == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.Major == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.Year < 1 || input.Year > 6 {
		return nil, appErrors.ErrInvalidInput
	}

	// Verify user exists and has the correct role
	u, err := s.userRepo.GetByID(input.UserID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}
	if u.Role != user.RoleStudent {
		return nil, appErrors.ErrRoleMismatch
	}

	// Check if student code already exists
	existing, _ := s.repo.GetByStudentCode(ctx, input.StudentCode)
	if existing != nil {
		return nil, appErrors.ErrAlreadyExists
	}

	// Check if user already has a student profile
	existing, _ = s.repo.GetByUserID(ctx, input.UserID)
	if existing != nil {
		return nil, appErrors.ErrAlreadyExists
	}

	st := &student.Student{
		UserID:      input.UserID,
		StudentCode: input.StudentCode,
		Major:       input.Major,
		Year:        input.Year,
		Status:      "enrolled",
		AdmittedAt:  input.AdmittedAt,
	}

	if err := s.repo.Create(ctx, st); err != nil {
		return nil, err
	}

	return st, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*student.Student, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	st, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrStudentNotFound
	}

	return st, nil
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (*student.Student, error) {
	if userID == "" {
		return nil, appErrors.ErrInvalidID
	}

	st, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrStudentNotFound
	}

	return st, nil
}

func (s *Service) GetWithDetails(ctx context.Context, id string) (*student.StudentWithDetails, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	details, err := s.repo.GetWithDetails(ctx, id)
	if err != nil {
		return nil, appErrors.ErrStudentNotFound
	}

	return details, nil
}

func (s *Service) List(ctx context.Context, filter *student.StudentFilter) ([]*student.Student, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	return s.repo.List(ctx, filter)
}

func (s *Service) UpdateStudent(ctx context.Context, id string, input *student.UpdateStudentInput) (*student.Student, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	st, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrStudentNotFound
	}

	if input.Major != nil {
		st.Major = *input.Major
	}
	if input.Year != nil {
		if *input.Year < 1 || *input.Year > 6 {
			return nil, appErrors.ErrInvalidInput
		}
		st.Year = *input.Year
	}
	if input.GPA != nil {
		if *input.GPA < 0 || *input.GPA > 4.0 {
			return nil, appErrors.ErrInvalidInput
		}
		st.GPA = *input.GPA
	}
	if input.Status != nil {
		validStatuses := map[string]bool{"enrolled": true, "on_leave": true, "graduated": true, "dropped": true}
		if !validStatuses[*input.Status] {
			return nil, appErrors.ErrInvalidInput
		}
		st.Status = *input.Status
	}

	if err := s.repo.Update(ctx, st); err != nil {
		return nil, err
	}

	return st, nil
}

func (s *Service) DeleteStudent(ctx context.Context, id string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrStudentNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) GetEnrollments(ctx context.Context, studentID string) ([]*student.StudentEnrollment, error) {
	if studentID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetEnrollments(ctx, studentID)
}

func (s *Service) GetGroups(ctx context.Context, studentID string) ([]*student.StudentGroup, error) {
	if studentID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetGroups(ctx, studentID)
}
