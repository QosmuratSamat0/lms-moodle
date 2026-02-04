package assignment

import (
	"context"
	"errors"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the assignment service interface
type Service interface {
	Create(ctx context.Context, teacherID uuid.UUID, req *CreateAssignmentRequest) (*AssignmentResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*AssignmentResponse, error)
	Update(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string, req *UpdateAssignmentRequest) error
	Delete(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string) error
	ListByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) (*AssignmentListResponse, error)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID, page, limit int) (*AssignmentListResponse, error)
	ListUpcoming(ctx context.Context, studentID uuid.UUID, limit int) ([]AssignmentResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new assignment service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, teacherID uuid.UUID, req *CreateAssignmentRequest) (*AssignmentResponse, error) {
	assignment := &Assignment{
		CourseID:           req.CourseID,
		GroupID:            req.GroupID,
		Title:              req.Title,
		Description:        req.Description,
		DueAt:              req.DueAt,
		MaxPoints:          req.MaxPoints,
		AllowLate:          req.AllowLate,
		CreatedByTeacherID: &teacherID,
	}

	if assignment.MaxPoints == 0 {
		assignment.MaxPoints = 100 // Default max points
	}

	if err := s.repo.Create(ctx, assignment); err != nil {
		return nil, errorx.Wrap(err, "create assignment")
	}

	// Fetch full details
	detailed, err := s.repo.GetByID(ctx, assignment.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created assignment")
	}

	return s.toResponse(detailed), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*AssignmentResponse, error) {
	assignment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("assignment")
		}
		return nil, errorx.Wrap(err, "get assignment")
	}
	return s.toResponse(assignment), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string, req *UpdateAssignmentRequest) error {
	assignment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("assignment")
		}
		return errorx.Wrap(err, "get assignment")
	}

	// Check ownership (admin can update any)
	if role != "admin" && (assignment.CreatedByTeacherID == nil || *assignment.CreatedByTeacherID != teacherID) {
		return errorx.NewForbiddenError("you don't have permission to update this assignment")
	}

	// Apply updates
	if req.Title != nil {
		assignment.Title = *req.Title
	}
	if req.Description != nil {
		assignment.Description = req.Description
	}
	if req.DueAt != nil {
		assignment.DueAt = req.DueAt
	}
	if req.MaxPoints != nil {
		assignment.MaxPoints = *req.MaxPoints
	}
	if req.AllowLate != nil {
		assignment.AllowLate = *req.AllowLate
	}
	if req.GroupID != nil {
		assignment.GroupID = req.GroupID
	}

	if err := s.repo.Update(ctx, &assignment.Assignment); err != nil {
		return errorx.Wrap(err, "update assignment")
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string) error {
	assignment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("assignment")
		}
		return errorx.Wrap(err, "get assignment")
	}

	// Check ownership (admin can delete any)
	if role != "admin" && (assignment.CreatedByTeacherID == nil || *assignment.CreatedByTeacherID != teacherID) {
		return errorx.NewForbiddenError("you don't have permission to delete this assignment")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errorx.Wrap(err, "delete assignment")
	}

	return nil
}

func (s *service) ListByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) (*AssignmentListResponse, error) {
	offset := (page - 1) * limit
	assignments, total, err := s.repo.ListByCourse(ctx, courseID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list course assignments")
	}

	return s.toListResponse(assignments, total, page, limit), nil
}

func (s *service) ListByTeacher(ctx context.Context, teacherID uuid.UUID, page, limit int) (*AssignmentListResponse, error) {
	offset := (page - 1) * limit
	assignments, total, err := s.repo.ListByTeacher(ctx, teacherID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list teacher assignments")
	}

	return s.toListResponse(assignments, total, page, limit), nil
}

func (s *service) ListUpcoming(ctx context.Context, studentID uuid.UUID, limit int) ([]AssignmentResponse, error) {
	assignments, err := s.repo.ListUpcoming(ctx, studentID, limit)
	if err != nil {
		return nil, errorx.Wrap(err, "list upcoming assignments")
	}

	var responses []AssignmentResponse
	for _, a := range assignments {
		responses = append(responses, *s.toResponse(&a))
	}
	return responses, nil
}

// Helper methods

func (s *service) toResponse(a *AssignmentWithDetails) *AssignmentResponse {
	resp := &AssignmentResponse{
		ID:               a.ID.String(),
		CourseID:         a.CourseID.String(),
		Title:            a.Title,
		Description:      a.Description,
		MaxPoints:        a.MaxPoints,
		AllowLate:        a.AllowLate,
		CourseTitle:      a.CourseTitle,
		GroupCode:        a.GroupCode,
		TeacherFirstName: a.TeacherFirstName,
		TeacherLastName:  a.TeacherLastName,
		SubmissionCount:  a.SubmissionCount,
		CreatedAt:        a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if a.DueAt != nil {
		dueStr := a.DueAt.Format("2006-01-02T15:04:05Z07:00")
		resp.DueAt = &dueStr
	}
	if a.CreatedByTeacherID != nil {
		id := a.CreatedByTeacherID.String()
		resp.CreatedByTeacher = &id
	}
	if a.GroupID != nil {
		id := a.GroupID.String()
		resp.GroupID = &id
	}
	return resp
}

func (s *service) toListResponse(assignments []AssignmentWithDetails, total int64, page, limit int) *AssignmentListResponse {
	var responses []AssignmentResponse
	for _, a := range assignments {
		responses = append(responses, *s.toResponse(&a))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &AssignmentListResponse{
		Assignments: responses,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
	}
}

var _ Service = (*service)(nil)
