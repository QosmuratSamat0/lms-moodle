package manager

import (
	"context"
	"errors"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the manager service interface
type Service interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*ManagerResponse, error)
	List(ctx context.Context, filter *ManagerFilter, page, limit int) (*ManagerListResponse, error)
	Update(ctx context.Context, userID uuid.UUID, req *UpdateManagerRequest) error
	GetSystemOverview(ctx context.Context) (*SystemOverview, error)
	ActivateUser(ctx context.Context, userID uuid.UUID) error
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
	BulkUserAction(ctx context.Context, req *BulkUserActionRequest) error
}

type service struct {
	repo Repository
}

// NewService creates a new manager service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, userID uuid.UUID) (*ManagerResponse, error) {
	manager, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("manager")
		}
		return nil, errorx.Wrap(err, "get manager")
	}

	return toManagerResponse(manager), nil
}

func (s *service) List(ctx context.Context, filter *ManagerFilter, page, limit int) (*ManagerListResponse, error) {
	managers, total, err := s.repo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, errorx.Wrap(err, "list managers")
	}

	responses := make([]ManagerResponse, len(managers))
	for i, manager := range managers {
		responses[i] = *toManagerResponse(&manager)
	}

	return &ManagerListResponse{
		Managers:   responses,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, req *UpdateManagerRequest) error {
	// Verify manager exists
	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("manager")
		}
		return errorx.Wrap(err, "get manager")
	}

	if err := s.repo.Update(ctx, userID, req); err != nil {
		return errorx.Wrap(err, "update manager")
	}

	return nil
}

func (s *service) GetSystemOverview(ctx context.Context) (*SystemOverview, error) {
	overview, err := s.repo.GetSystemOverview(ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "get system overview")
	}

	return overview, nil
}

func (s *service) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.ActivateUser(ctx, userID); err != nil {
		return errorx.Wrap(err, "activate user")
	}
	return nil
}

func (s *service) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.DeactivateUser(ctx, userID); err != nil {
		return errorx.Wrap(err, "deactivate user")
	}
	return nil
}

func (s *service) BulkUserAction(ctx context.Context, req *BulkUserActionRequest) error {
	switch req.Action {
	case "activate":
		if err := s.repo.BulkActivateUsers(ctx, req.UserIDs); err != nil {
			return errorx.Wrap(err, "bulk activate users")
		}
	case "deactivate":
		if err := s.repo.BulkDeactivateUsers(ctx, req.UserIDs); err != nil {
			return errorx.Wrap(err, "bulk deactivate users")
		}
	default:
		return errorx.NewValidationError("invalid action")
	}
	return nil
}

func toManagerResponse(m *ManagerWithUser) *ManagerResponse {
	return &ManagerResponse{
		UserID:    m.UserID,
		Email:     m.Email,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		IsActive:  m.IsActive,
	}
}
