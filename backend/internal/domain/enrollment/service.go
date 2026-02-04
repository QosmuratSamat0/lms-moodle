package enrollment

import (
	"context"
	"errors"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the enrollment service interface
type Service interface {
	Enroll(ctx context.Context, studentID, courseID uuid.UUID) (*EnrollmentResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*EnrollmentResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status EnrollmentStatus, requesterID uuid.UUID, role string) error
	Drop(ctx context.Context, enrollmentID, studentID uuid.UUID) error
	ListByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) (*EnrollmentListResponse, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, page, limit int) (*EnrollmentListResponse, error)
	IsEnrolled(ctx context.Context, courseID, studentID uuid.UUID) (bool, error)
}

type service struct {
	repo Repository
}

// NewService creates a new enrollment service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Enroll(ctx context.Context, studentID, courseID uuid.UUID) (*EnrollmentResponse, error) {
	// Check if already enrolled
	existing, err := s.repo.GetByCourseAndStudent(ctx, courseID, studentID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.Wrap(err, "check existing enrollment")
	}
	if existing != nil {
		if existing.Status == StatusActive {
			return nil, errorx.NewConflictError("already enrolled in this course")
		}
		if existing.Status == StatusPending {
			return nil, errorx.NewConflictError("enrollment request already pending")
		}
	}

	enrollment := &Enrollment{
		CourseID:  courseID,
		StudentID: studentID,
		Status:    StatusActive, // Auto-approve for now; can be changed to StatusPending
	}

	if err := s.repo.Create(ctx, enrollment); err != nil {
		return nil, errorx.Wrap(err, "create enrollment")
	}

	// Fetch full details
	detailed, err := s.repo.GetByID(ctx, enrollment.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get enrollment details")
	}

	return s.toResponse(detailed), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*EnrollmentResponse, error) {
	enrollment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("enrollment")
		}
		return nil, errorx.Wrap(err, "get enrollment")
	}
	return s.toResponse(enrollment), nil
}

func (s *service) UpdateStatus(ctx context.Context, id uuid.UUID, status EnrollmentStatus, requesterID uuid.UUID, role string) error {
	// Only teachers/admins can update enrollment status
	if role != "teacher" && role != "admin" && role != "manager" {
		return errorx.NewForbiddenError("only teachers, managers, or admins can update enrollment status")
	}

	enrollment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("enrollment")
		}
		return errorx.Wrap(err, "get enrollment")
	}

	// For teachers, verify they own the course (simplified check - in real app, verify via course ownership)
	_ = enrollment // Could add more ownership checks here

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return errorx.Wrap(err, "update enrollment status")
	}

	return nil
}

func (s *service) Drop(ctx context.Context, enrollmentID, studentID uuid.UUID) error {
	enrollment, err := s.repo.GetByID(ctx, enrollmentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("enrollment")
		}
		return errorx.Wrap(err, "get enrollment")
	}

	// Verify the student owns this enrollment
	if enrollment.StudentID != studentID {
		return errorx.NewForbiddenError("you can only drop your own enrollments")
	}

	if err := s.repo.UpdateStatus(ctx, enrollmentID, StatusDropped); err != nil {
		return errorx.Wrap(err, "drop enrollment")
	}

	return nil
}

func (s *service) ListByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) (*EnrollmentListResponse, error) {
	offset := (page - 1) * limit
	enrollments, total, err := s.repo.ListByCourse(ctx, courseID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list course enrollments")
	}

	return s.toListResponse(enrollments, total, page, limit), nil
}

func (s *service) ListByStudent(ctx context.Context, studentID uuid.UUID, page, limit int) (*EnrollmentListResponse, error) {
	offset := (page - 1) * limit
	enrollments, total, err := s.repo.ListByStudent(ctx, studentID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list student enrollments")
	}

	return s.toListResponse(enrollments, total, page, limit), nil
}

func (s *service) IsEnrolled(ctx context.Context, courseID, studentID uuid.UUID) (bool, error) {
	return s.repo.IsEnrolled(ctx, courseID, studentID)
}

// Helper methods

func (s *service) toResponse(e *EnrollmentWithDetails) *EnrollmentResponse {
	return &EnrollmentResponse{
		ID:               e.ID.String(),
		CourseID:         e.CourseID.String(),
		StudentID:        e.StudentID.String(),
		Status:           string(e.Status),
		CourseTitle:      e.CourseTitle,
		StudentFirstName: e.StudentFirstName,
		StudentLastName:  e.StudentLastName,
		StudentEmail:     e.StudentEmail,
		EnrolledAt:       e.EnrolledAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *service) toListResponse(enrollments []EnrollmentWithDetails, total int64, page, limit int) *EnrollmentListResponse {
	var responses []EnrollmentResponse
	for _, e := range enrollments {
		responses = append(responses, *s.toResponse(&e))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &EnrollmentListResponse{
		Enrollments: responses,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
	}
}

var _ Service = (*service)(nil)
