package coursecategory

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/coursecategory"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo coursecategory.Repository
}

func NewService(repo coursecategory.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input *coursecategory.CreateCourseCategoryInput) (*coursecategory.CourseCategory, error) {
	if input.Name == "" {
		return nil, appErrors.ErrMissingRequired
	}

	cc := &coursecategory.CourseCategory{
		Name:        input.Name,
		Description: input.Description,
		Icon:        input.Icon,
		Order:       input.Order,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, cc); err != nil {
		return nil, err
	}

	return cc, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*coursecategory.CourseCategory, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	cc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrCategoryNotFound
	}

	return cc, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*coursecategory.CourseCategory, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filter := &coursecategory.CourseCategoryFilter{
		Limit:  limit,
		Offset: offset,
	}

	categories, _, err := s.repo.List(ctx, filter)
	return categories, err
}

func (s *Service) Update(ctx context.Context, id string, input *coursecategory.UpdateCourseCategoryInput) (*coursecategory.CourseCategory, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	cc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrCategoryNotFound
	}

	if input.Name != nil {
		cc.Name = *input.Name
	}
	if input.Description != nil {
		cc.Description = *input.Description
	}
	if input.Icon != nil {
		cc.Icon = *input.Icon
	}
	if input.Order != nil {
		cc.Order = *input.Order
	}
	if input.IsActive != nil {
		cc.IsActive = *input.IsActive
	}

	if err := s.repo.Update(ctx, cc); err != nil {
		return nil, err
	}

	return cc, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrCategoryNotFound
	}

	return s.repo.Delete(ctx, id)
}
