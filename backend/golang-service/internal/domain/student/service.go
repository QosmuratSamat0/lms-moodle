package student

import (
	"context"
	"errors"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the student service interface
type Service interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*StudentResponse, error)
	List(ctx context.Context, filter *StudentFilter, page, limit int) (*StudentListResponse, error)
	Update(ctx context.Context, userID uuid.UUID, req *UpdateStudentRequest) error
	GetStats(ctx context.Context, userID uuid.UUID) (*StudentStats, error)
	ListByGroup(ctx context.Context, groupName string) ([]StudentResponse, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]StudentResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new student service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, userID uuid.UUID) (*StudentResponse, error) {
	student, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("student")
		}
		return nil, errorx.Wrap(err, "get student")
	}

	return toStudentResponse(student), nil
}

func (s *service) List(ctx context.Context, filter *StudentFilter, page, limit int) (*StudentListResponse, error) {
	students, total, err := s.repo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, errorx.Wrap(err, "list students")
	}

	responses := make([]StudentResponse, len(students))
	for i, student := range students {
		responses[i] = *toStudentResponse(&student)
	}

	return &StudentListResponse{
		Students:   responses,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, req *UpdateStudentRequest) error {
	// Verify student exists
	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("student")
		}
		return errorx.Wrap(err, "get student")
	}

	if err := s.repo.Update(ctx, userID, req); err != nil {
		return errorx.Wrap(err, "update student")
	}

	return nil
}

func (s *service) GetStats(ctx context.Context, userID uuid.UUID) (*StudentStats, error) {
	stats, err := s.repo.GetStats(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("student")
		}
		return nil, errorx.Wrap(err, "get student stats")
	}

	return stats, nil
}

func (s *service) ListByGroup(ctx context.Context, groupName string) ([]StudentResponse, error) {
	students, err := s.repo.ListByGroup(ctx, groupName)
	if err != nil {
		return nil, errorx.Wrap(err, "list students by group")
	}

	responses := make([]StudentResponse, len(students))
	for i, student := range students {
		responses[i] = *toStudentResponse(&student)
	}

	return responses, nil
}

func (s *service) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]StudentResponse, error) {
	students, err := s.repo.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, errorx.Wrap(err, "list students by course")
	}

	responses := make([]StudentResponse, len(students))
	for i, student := range students {
		responses[i] = *toStudentResponse(&student)
	}

	return responses, nil
}

func toStudentResponse(s *StudentWithUser) *StudentResponse {
	return &StudentResponse{
		UserID:    s.UserID,
		Email:     s.Email,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		GroupName: s.GroupName,
		IsActive:  s.IsActive,
	}
}
