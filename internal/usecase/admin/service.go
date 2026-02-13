package admin

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/admin"
	"github.com/ap1-final-mini-moodle/internal/domain/user"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo     admin.Repository
	userRepo user.Repository
}

func NewService(repo admin.Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) CreateAdmin(ctx context.Context, input *admin.CreateAdminInput) (*admin.Admin, error) {
	if input.UserID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.EmployeeID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.AccessLevel == "" {
		return nil, appErrors.ErrMissingRequired
	}

	// Validate access level
	if !isValidAccessLevel(input.AccessLevel) {
		return nil, appErrors.ErrInvalidInput
	}

	// Verify user exists and has the correct role
	u, err := s.userRepo.GetByID(input.UserID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}
	if u.Role != user.RoleAdmin {
		return nil, appErrors.ErrRoleMismatch
	}

	// Check if employee ID already exists
	existing, err := s.repo.GetByEmployeeID(ctx, input.EmployeeID)
	if err != nil && err != appErrors.ErrAdminNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, appErrors.ErrAlreadyExists
	}

	// Check if user already has an admin profile
	existing, err = s.repo.GetByUserID(ctx, input.UserID)
	if err != nil && err != appErrors.ErrAdminNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, appErrors.ErrAlreadyExists
	}

	a := &admin.Admin{
		UserID:      input.UserID,
		EmployeeID:  input.EmployeeID,
		Department:  input.Department,
		AccessLevel: input.AccessLevel,
		Permissions: input.Permissions,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

func isValidAccessLevel(level string) bool {
	switch level {
	case "admin", "super_admin":
		return true
	default:
		return false
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*admin.Admin, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrAdminNotFound
	}

	return a, nil
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (*admin.Admin, error) {
	if userID == "" {
		return nil, appErrors.ErrInvalidID
	}

	a, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrAdminNotFound
	}

	return a, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*admin.Admin, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filter := &admin.AdminFilter{
		Limit:  limit,
		Offset: offset,
	}

	admins, _, err := s.repo.List(ctx, filter)
	return admins, err
}

func (s *Service) UpdateAdmin(ctx context.Context, id string, input *admin.UpdateAdminInput) (*admin.Admin, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrAdminNotFound
	}

	if input.Department != nil {
		a.Department = *input.Department
	}
	if input.AccessLevel != nil {
		// Validate access level
		if *input.AccessLevel != "super_admin" && *input.AccessLevel != "admin" {
			return nil, appErrors.ErrInvalidInput
		}
		a.AccessLevel = *input.AccessLevel
	}
	if input.Permissions != nil {
		a.Permissions = *input.Permissions
	}
	if input.IsActive != nil {
		a.IsActive = *input.IsActive
	}

	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) DeleteAdmin(ctx context.Context, id string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrAdminNotFound
	}

	return s.repo.Delete(ctx, id)
}
