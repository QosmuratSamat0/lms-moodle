package categorymanager

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/categorymanager"
	"github.com/ap1-final-mini-moodle/internal/domain/user"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo     categorymanager.Repository
	userRepo user.Repository
}

func NewService(repo categorymanager.Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) Create(ctx context.Context, input *categorymanager.CreateCategoryManagerInput) (*categorymanager.CategoryManager, error) {
	if input.UserID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.CategoryID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.PermissionLevel == "" {
		return nil, appErrors.ErrMissingRequired
	}

	// Validate user role is manager
	u, err := s.userRepo.GetByID(input.UserID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}
	if u.Role != user.RoleManager {
		return nil, appErrors.ErrRoleMismatch
	}

	cm := &categorymanager.CategoryManager{
		UserID:          input.UserID,
		CategoryID:      input.CategoryID,
		PermissionLevel: input.PermissionLevel,
		IsActive:        true,
	}

	if err := s.repo.Create(ctx, cm); err != nil {
		return nil, err
	}

	return cm, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*categorymanager.CategoryManager, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	cm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrCategoryManagerNotFound
	}

	return cm, nil
}

func (s *Service) GetByUserAndCategory(ctx context.Context, userID, categoryID string) (*categorymanager.CategoryManager, error) {
	if userID == "" || categoryID == "" {
		return nil, appErrors.ErrInvalidID
	}

	cm, err := s.repo.GetByUserAndCategory(ctx, userID, categoryID)
	if err != nil {
		return nil, appErrors.ErrCategoryManagerNotFound
	}

	return cm, nil
}

func (s *Service) GetByUserID(ctx context.Context, userID string) ([]*categorymanager.CategoryManager, error) {
	if userID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) GetByCategoryID(ctx context.Context, categoryID string) ([]*categorymanager.CategoryManager, error) {
	if categoryID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetByCategoryID(ctx, categoryID)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*categorymanager.CategoryManager, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filter := &categorymanager.CategoryManagerFilter{
		Limit:  limit,
		Offset: offset,
	}

	managers, _, err := s.repo.List(ctx, filter)
	return managers, err
}

func (s *Service) Update(ctx context.Context, id string, input *categorymanager.UpdateCategoryManagerInput) (*categorymanager.CategoryManager, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	cm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrCategoryManagerNotFound
	}

	if input.PermissionLevel != nil {
		// Validate permission level
		if *input.PermissionLevel != "view" && *input.PermissionLevel != "edit" && *input.PermissionLevel != "admin" {
			return nil, appErrors.ErrInvalidInput
		}
		cm.PermissionLevel = *input.PermissionLevel
	}
	if input.IsActive != nil {
		cm.IsActive = *input.IsActive
	}

	if err := s.repo.Update(ctx, cm); err != nil {
		return nil, err
	}

	return cm, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrCategoryManagerNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) GetWithDetails(ctx context.Context, id string) (*categorymanager.CategoryManagerWithDetails, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetWithDetails(ctx, id)
}
