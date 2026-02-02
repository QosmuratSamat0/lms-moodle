package teacher

import (
	"context"
	"errors"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the teacher service interface
type Service interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*TeacherResponse, error)
	List(ctx context.Context, filter *TeacherFilter, page, limit int) (*TeacherListResponse, error)
	Update(ctx context.Context, userID uuid.UUID, req *UpdateTeacherRequest) error
	GetStats(ctx context.Context, userID uuid.UUID) (*TeacherStats, error)
	ListByDepartment(ctx context.Context, department string) ([]TeacherResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new teacher service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, userID uuid.UUID) (*TeacherResponse, error) {
	teacher, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("teacher")
		}
		return nil, errorx.Wrap(err, "get teacher")
	}

	return toTeacherResponse(teacher), nil
}

func (s *service) List(ctx context.Context, filter *TeacherFilter, page, limit int) (*TeacherListResponse, error) {
	teachers, total, err := s.repo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, errorx.Wrap(err, "list teachers")
	}

	responses := make([]TeacherResponse, len(teachers))
	for i, teacher := range teachers {
		responses[i] = *toTeacherResponse(&teacher)
	}

	return &TeacherListResponse{
		Teachers:   responses,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, req *UpdateTeacherRequest) error {
	// Verify teacher exists
	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("teacher")
		}
		return errorx.Wrap(err, "get teacher")
	}

	if err := s.repo.Update(ctx, userID, req); err != nil {
		return errorx.Wrap(err, "update teacher")
	}

	return nil
}

func (s *service) GetStats(ctx context.Context, userID uuid.UUID) (*TeacherStats, error) {
	stats, err := s.repo.GetStats(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("teacher")
		}
		return nil, errorx.Wrap(err, "get teacher stats")
	}

	return stats, nil
}

func (s *service) ListByDepartment(ctx context.Context, department string) ([]TeacherResponse, error) {
	teachers, err := s.repo.ListByDepartment(ctx, department)
	if err != nil {
		return nil, errorx.Wrap(err, "list teachers by department")
	}

	responses := make([]TeacherResponse, len(teachers))
	for i, teacher := range teachers {
		responses[i] = *toTeacherResponse(&teacher)
	}

	return responses, nil
}

func toTeacherResponse(t *TeacherWithUser) *TeacherResponse {
	return &TeacherResponse{
		UserID:     t.UserID,
		Email:      t.Email,
		FirstName:  t.FirstName,
		LastName:   t.LastName,
		Department: t.Department,
		IsActive:   t.IsActive,
	}
}
