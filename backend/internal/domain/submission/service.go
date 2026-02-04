package submission

import (
	"context"
	"errors"
	"time"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the submission service interface
type Service interface {
	Submit(ctx context.Context, studentID uuid.UUID, req *CreateSubmissionRequest) (*SubmissionResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SubmissionResponse, error)
	Update(ctx context.Context, id uuid.UUID, studentID uuid.UUID, req *UpdateSubmissionRequest) error
	Delete(ctx context.Context, id uuid.UUID, studentID uuid.UUID, role string) error
	ListByAssignment(ctx context.Context, assignmentID uuid.UUID, groupID *uuid.UUID, page, limit int) (*SubmissionListResponse, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, page, limit int) (*SubmissionListResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new submission service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Submit(ctx context.Context, studentID uuid.UUID, req *CreateSubmissionRequest) (*SubmissionResponse, error) {
	// First, get assignment details to check deadline and allow_late
	assignmentInfo, err := s.repo.GetAssignmentInfo(ctx, req.AssignmentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("assignment")
		}
		return nil, errorx.Wrap(err, "get assignment info")
	}

	// Check if submission is past deadline and late submissions are not allowed
	now := time.Now()
	isLate := assignmentInfo.DueAt != nil && now.After(*assignmentInfo.DueAt)
	if isLate && !assignmentInfo.AllowLate {
		return nil, errorx.NewBadRequestError("submission deadline has passed and late submissions are not allowed")
	}

	// Check if already submitted
	existing, err := s.repo.GetByAssignmentAndStudent(ctx, req.AssignmentID, studentID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.Wrap(err, "check existing submission")
	}

	if existing != nil {
		// Update existing submission (resubmit)
		existing.ContentText = req.ContentText
		existing.FileURL = req.FileURL
		if isLate {
			existing.Status = StatusLate
		} else {
			existing.Status = StatusResubmitted
		}

		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, errorx.Wrap(err, "update submission")
		}

		detailed, err := s.repo.GetByID(ctx, existing.ID)
		if err != nil {
			return nil, errorx.Wrap(err, "get updated submission")
		}
		return s.toResponse(detailed), nil
	}

	// Determine initial status
	status := StatusSubmitted
	if isLate {
		status = StatusLate
	}

	// Create new submission
	submission := &Submission{
		AssignmentID: req.AssignmentID,
		StudentID:    studentID,
		ContentText:  req.ContentText,
		FileURL:      req.FileURL,
		Status:       status,
	}

	if err := s.repo.Create(ctx, submission); err != nil {
		return nil, errorx.Wrap(err, "create submission")
	}

	// Fetch full details
	detailed, err := s.repo.GetByID(ctx, submission.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created submission")
	}

	return s.toResponse(detailed), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*SubmissionResponse, error) {
	submission, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("submission")
		}
		return nil, errorx.Wrap(err, "get submission")
	}
	return s.toResponse(submission), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, studentID uuid.UUID, req *UpdateSubmissionRequest) error {
	submission, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("submission")
		}
		return errorx.Wrap(err, "get submission")
	}

	// Verify ownership
	if submission.StudentID != studentID {
		return errorx.NewForbiddenError("you can only update your own submissions")
	}

	// Can't update if already graded
	if submission.GradeScore != nil {
		return errorx.NewBadRequestError("cannot update a graded submission")
	}

	// Apply updates
	if req.ContentText != nil {
		submission.ContentText = req.ContentText
	}
	if req.FileURL != nil {
		submission.FileURL = req.FileURL
	}
	submission.Status = StatusResubmitted

	if err := s.repo.Update(ctx, &submission.Submission); err != nil {
		return errorx.Wrap(err, "update submission")
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, studentID uuid.UUID, role string) error {
	submission, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("submission")
		}
		return errorx.Wrap(err, "get submission")
	}

	// Verify ownership (admin/teacher can delete any)
	if role != "admin" && role != "teacher" && submission.StudentID != studentID {
		return errorx.NewForbiddenError("you can only delete your own submissions")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errorx.Wrap(err, "delete submission")
	}

	return nil
}

func (s *service) ListByAssignment(ctx context.Context, assignmentID uuid.UUID, groupID *uuid.UUID, page, limit int) (*SubmissionListResponse, error) {
	offset := (page - 1) * limit
	submissions, total, err := s.repo.ListByAssignment(ctx, assignmentID, groupID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list assignment submissions")
	}

	return s.toListResponse(submissions, total, page, limit), nil
}

func (s *service) ListByStudent(ctx context.Context, studentID uuid.UUID, page, limit int) (*SubmissionListResponse, error) {
	offset := (page - 1) * limit
	submissions, total, err := s.repo.ListByStudent(ctx, studentID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list student submissions")
	}

	return s.toListResponse(submissions, total, page, limit), nil
}

// Helper methods

func (s *service) toResponse(sub *SubmissionWithDetails) *SubmissionResponse {
	resp := &SubmissionResponse{
		ID:               sub.ID.String(),
		AssignmentID:     sub.AssignmentID.String(),
		StudentID:        sub.StudentID.String(),
		ContentText:      sub.ContentText,
		FileURL:          sub.FileURL,
		Status:           string(sub.Status),
		SubmittedAt:      sub.SubmittedAt.Format("2006-01-02T15:04:05Z07:00"),
		AssignmentTitle:  sub.AssignmentTitle,
		CourseTitle:      sub.CourseTitle,
		StudentFirstName: sub.StudentFirstName,
		StudentLastName:  sub.StudentLastName,
		StudentEmail:     sub.StudentEmail,
		MaxPoints:        sub.MaxPoints,
		GradeScore:       sub.GradeScore,
		GradeFeedback:    sub.GradeFeedback,
	}
	if sub.DueAt != nil {
		dueStr := sub.DueAt.Format("2006-01-02T15:04:05Z07:00")
		resp.DueAt = &dueStr
	}
	return resp
}

func (s *service) toListResponse(submissions []SubmissionWithDetails, total int64, page, limit int) *SubmissionListResponse {
	responses := make([]SubmissionResponse, 0, len(submissions))
	for _, sub := range submissions {
		responses = append(responses, *s.toResponse(&sub))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &SubmissionListResponse{
		Submissions: responses,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
	}
}

var _ Service = (*service)(nil)
