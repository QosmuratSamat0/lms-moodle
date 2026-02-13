package manager

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/manager"
	"github.com/ap1-final-mini-moodle/internal/domain/user"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo     manager.Repository
	userRepo user.Repository
}

func NewService(repo manager.Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) CreateManager(ctx context.Context, input *manager.CreateManagerInput) (*manager.Manager, error) {
	if input.UserID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.EmployeeID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.Department == "" {
		return nil, appErrors.ErrMissingRequired
	}

	// Verify user exists and has the correct role
	u, err := s.userRepo.GetByID(input.UserID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}
	if u.Role != user.RoleManager {
		return nil, appErrors.ErrRoleMismatch
	}

	// Check if employee ID already exists
	existing, _ := s.repo.GetByEmployeeID(ctx, input.EmployeeID)
	if existing != nil {
		return nil, appErrors.ErrAlreadyExists
	}

	// Check if user already has a manager profile
	existing, _ = s.repo.GetByUserID(ctx, input.UserID)
	if existing != nil {
		return nil, appErrors.ErrAlreadyExists
	}

	m := &manager.Manager{
		UserID:            input.UserID,
		EmployeeID:        input.EmployeeID,
		Department:        input.Department,
		ManagesCategories: input.ManagesCategories,
		ManagesTeachers:   input.ManagesTeachers,
		IsActive:          true,
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*manager.Manager, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrManagerNotFound
	}

	return m, nil
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (*manager.Manager, error) {
	if userID == "" {
		return nil, appErrors.ErrInvalidID
	}

	m, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrManagerNotFound
	}

	return m, nil
}

func (s *Service) GetByDepartment(ctx context.Context, department string) ([]*manager.Manager, error) {
	if department == "" {
		return nil, appErrors.ErrInvalidInput
	}

	return s.repo.GetByDepartment(ctx, department)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*manager.Manager, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filter := &manager.ManagerFilter{
		Limit:  limit,
		Offset: offset,
	}

	managers, _, err := s.repo.List(ctx, filter)
	return managers, err
}

func (s *Service) UpdateManager(ctx context.Context, id string, input *manager.UpdateManagerInput) (*manager.Manager, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrManagerNotFound
	}

	if input.Department != nil {
		m.Department = *input.Department
	}
	if input.ManagesCategories != nil {
		m.ManagesCategories = *input.ManagesCategories
	}
	if input.ManagesTeachers != nil {
		m.ManagesTeachers = *input.ManagesTeachers
	}
	if input.IsActive != nil {
		m.IsActive = *input.IsActive
	}

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Service) DeleteManager(ctx context.Context, id string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrManagerNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) AddManagedCategory(ctx context.Context, managerID, categoryID string) error {
	if managerID == "" || categoryID == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, managerID)
	if err != nil {
		return appErrors.ErrManagerNotFound
	}

	return s.repo.AddManagedCategory(ctx, managerID, categoryID)
}

func (s *Service) RemoveManagedCategory(ctx context.Context, managerID, categoryID string) error {
	if managerID == "" || categoryID == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, managerID)
	if err != nil {
		return appErrors.ErrManagerNotFound
	}

	return s.repo.RemoveManagedCategory(ctx, managerID, categoryID)
}

func (s *Service) AddManagedTeacher(ctx context.Context, managerID, teacherID string) error {
	if managerID == "" || teacherID == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, managerID)
	if err != nil {
		return appErrors.ErrManagerNotFound
	}

	return s.repo.AddManagedTeacher(ctx, managerID, teacherID)
}

func (s *Service) RemoveManagedTeacher(ctx context.Context, managerID, teacherID string) error {
	if managerID == "" || teacherID == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, managerID)
	if err != nil {
		return appErrors.ErrManagerNotFound
	}

	return s.repo.RemoveManagedTeacher(ctx, managerID, teacherID)
}

func (s *Service) GetWithDetails(ctx context.Context, id string) (*manager.ManagerWithDetails, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetWithDetails(ctx, id)
}
