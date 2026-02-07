package appeal

import (
	"context"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/appeal"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo appeal.Repository
}

func NewService(repo appeal.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAppeal(ctx context.Context, studentID string, input *appeal.CreateAppealInput) (*appeal.GradeAppeal, error) {
	if input.GradeID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.Reason == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if len(input.Reason) < 10 {
		return nil, appErrors.ErrInvalidInput
	}

	// Check if there's already a pending appeal
	hasPending, _ := s.repo.HasPendingAppeal(ctx, input.GradeID)
	if hasPending {
		return nil, appErrors.ErrAppealAlreadyExists
	}

	a := &appeal.GradeAppeal{
		GradeID:   input.GradeID,
		StudentID: studentID,
		Reason:    input.Reason,
		Evidence:  input.Evidence,
		Status:    appeal.AppealPending,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) GetAppealByID(ctx context.Context, id string) (*appeal.GradeAppeal, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrAppealNotFound
	}

	return a, nil
}

func (s *Service) GetStudentAppeals(ctx context.Context, studentID string, status *appeal.AppealStatus, limit, offset int) ([]*appeal.AppealWithDetails, int64, error) {
	if studentID == "" {
		return nil, 0, appErrors.ErrInvalidID
	}

	if limit <= 0 {
		limit = 20
	}

	return s.repo.GetByStudentID(ctx, studentID, status, limit, offset)
}

func (s *Service) GetTeacherAppeals(ctx context.Context, teacherID string, status *appeal.AppealStatus, limit, offset int) ([]*appeal.AppealWithDetails, int64, error) {
	if teacherID == "" {
		return nil, 0, appErrors.ErrInvalidID
	}

	if limit <= 0 {
		limit = 20
	}

	return s.repo.GetByTeacherCourses(ctx, teacherID, status, limit, offset)
}

func (s *Service) ResolveAppeal(ctx context.Context, id string, teacherID string, input *appeal.ResolveAppealInput) (*appeal.GradeAppeal, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}
	if input.Response == "" {
		return nil, appErrors.ErrMissingRequired
	}

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrAppealNotFound
	}

	// Check if already resolved
	if a.Status != appeal.AppealPending {
		return nil, appErrors.ErrAppealAlreadyResolved
	}

	// Validate status
	if input.Status != appeal.AppealApproved && input.Status != appeal.AppealRejected {
		return nil, appErrors.ErrInvalidInput
	}

	// If approving, new score must be provided
	if input.Status == appeal.AppealApproved && input.NewScore == nil {
		return nil, appErrors.ErrMissingRequired
	}

	now := time.Now()
	a.Status = input.Status
	a.TeacherID = &teacherID
	a.Response = &input.Response
	a.NewScore = input.NewScore
	a.ResolvedAt = &now

	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Service) DeleteAppeal(ctx context.Context, id string, studentID string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrAppealNotFound
	}

	// Only allow deletion of own pending appeals
	if a.StudentID != studentID {
		return appErrors.ErrForbidden
	}
	if a.Status != appeal.AppealPending {
		return appErrors.ErrAppealAlreadyResolved
	}

	return s.repo.Delete(ctx, id)
}
