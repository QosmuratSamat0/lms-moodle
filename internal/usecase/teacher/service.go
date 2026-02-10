package teacher

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/teacher"
	"github.com/ap1-final-mini-moodle/internal/domain/user"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo     teacher.Repository
	userRepo user.Repository
}

func NewService(repo teacher.Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) CreateTeacher(ctx context.Context, input *teacher.CreateTeacherInput) (*teacher.Teacher, error) {
	if input.UserID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.EmployeeID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.FirstName == "" || input.LastName == "" {
		return nil, appErrors.ErrMissingRequired
	}

	// Verify user exists and has the correct role
	u, err := s.userRepo.GetByID(input.UserID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}
	if u.Role != user.RoleTeacher {
		return nil, appErrors.ErrRoleMismatch
	}

	// Check if employee ID already exists
	existing, _ := s.repo.GetByEmployeeID(ctx, input.EmployeeID)
	if existing != nil {
		return nil, appErrors.ErrAlreadyExists
	}

	// Check if user already has a teacher profile
	existing, _ = s.repo.GetByUserID(ctx, input.UserID)
	if existing != nil {
		return nil, appErrors.ErrAlreadyExists
	}

	t := &teacher.Teacher{
		UserID:         input.UserID,
		EmployeeID:     input.EmployeeID,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Department:     input.Department,
		Specialization: input.Specialization,
		Qualifications: input.Qualifications,
		Bio:            input.Bio,
		OfficeHours:    input.OfficeHours,
		Phone:          input.Phone,
		IsActive:       true,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*teacher.Teacher, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrTeacherNotFound
	}

	return t, nil
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (*teacher.Teacher, error) {
	if userID == "" {
		return nil, appErrors.ErrInvalidID
	}

	t, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrTeacherNotFound
	}

	return t, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*teacher.Teacher, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filter := &teacher.TeacherFilter{
		Limit:  limit,
		Offset: offset,
	}

	teachers, _, err := s.repo.List(ctx, filter)
	return teachers, err
}

func (s *Service) UpdateTeacher(ctx context.Context, id string, input *teacher.UpdateTeacherInput) (*teacher.Teacher, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrTeacherNotFound
	}

	if input.Department != nil {
		t.Department = *input.Department
	}
	if input.Specialization != nil {
		t.Specialization = *input.Specialization
	}
	if input.Qualifications != nil {
		t.Qualifications = *input.Qualifications
	}
	if input.Bio != nil {
		t.Bio = *input.Bio
	}
	if input.OfficeHours != nil {
		t.OfficeHours = *input.OfficeHours
	}
	if input.Phone != nil {
		t.Phone = *input.Phone
	}
	if input.IsActive != nil {
		t.IsActive = *input.IsActive
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) DeleteTeacher(ctx context.Context, id string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrTeacherNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) GetTeacherCourses(ctx context.Context, teacherID string) ([]teacher.TeacherCourse, error) {
	if teacherID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetTeacherCourses(ctx, teacherID)
}

func (s *Service) GetTeacherGroups(ctx context.Context, teacherID string) ([]*teacher.TeacherGroup, error) {
	if teacherID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetTeacherGroups(ctx, teacherID)
}
