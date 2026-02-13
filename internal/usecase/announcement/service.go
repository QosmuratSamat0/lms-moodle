package announcement

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/announcement"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo announcement.Repository
}

func NewService(repo announcement.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, authorID string, input *announcement.CreateAnnouncementInput) (*announcement.Announcement, error) {
	if input.CourseID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.Title == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.Content == "" {
		return nil, appErrors.ErrMissingRequired
	}

	a := &announcement.Announcement{
		CourseID: input.CourseID,
		AuthorID: authorID,
		Title:    input.Title,
		Content:  input.Content,
		Pinned:   input.Pinned,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*announcement.Announcement, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrAnnouncementNotFound
	}

	return a, nil
}

func (s *Service) GetByCourseID(ctx context.Context, courseID string, limit, offset int) ([]*announcement.Announcement, int64, error) {
	if courseID == "" {
		return nil, 0, appErrors.ErrInvalidID
	}

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetByCourseID(ctx, courseID, limit, offset)
}

func (s *Service) Update(ctx context.Context, id string, authorID string, input *announcement.UpdateAnnouncementInput) (*announcement.Announcement, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrAnnouncementNotFound
	}

	// Verify author
	if a.AuthorID != authorID {
		return nil, appErrors.ErrForbidden
	}

	if input.Title != nil {
		a.Title = *input.Title
	}
	if input.Content != nil {
		a.Content = *input.Content
	}
	if input.Pinned != nil {
		a.Pinned = *input.Pinned
	}

	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) Delete(ctx context.Context, id string, authorID string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrAnnouncementNotFound
	}

	// Verify author
	if a.AuthorID != authorID {
		return appErrors.ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}
