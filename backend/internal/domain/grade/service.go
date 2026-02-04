package grade

import (
	"context"
	"errors"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the grade service interface
type Service interface {
	GradeSubmission(ctx context.Context, teacherID uuid.UUID, req *GradeSubmissionRequest) (*GradeResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*GradeResponse, error)
	Update(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string, req *UpdateGradeRequest) error
	Delete(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string) error
	ListByStudent(ctx context.Context, studentID uuid.UUID, page, limit int) (*GradeListResponse, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) (*GradeListResponse, error)
	ListByAssignment(ctx context.Context, assignmentID uuid.UUID, page, limit int) (*GradeListResponse, error)
	GetStudentCourseSummary(ctx context.Context, studentID, courseID uuid.UUID) (*StudentGradeSummary, error)
}

type service struct {
	repo Repository
}

// NewService creates a new grade service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GradeSubmission(ctx context.Context, teacherID uuid.UUID, req *GradeSubmissionRequest) (*GradeResponse, error) {
	// Check if already graded
	existing, err := s.repo.GetBySubmissionID(ctx, req.SubmissionID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.Wrap(err, "check existing grade")
	}

	if existing != nil {
		// Update existing grade
		existing.Score = req.Score
		existing.Feedback = req.Feedback
		existing.GradedByTeacherID = &teacherID

		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, errorx.Wrap(err, "update existing grade")
		}

		detailed, err := s.repo.GetByID(ctx, existing.ID)
		if err != nil {
			return nil, errorx.Wrap(err, "get updated grade")
		}
		return s.toResponse(detailed), nil
	}

	// Create new grade
	grade := &Grade{
		SubmissionID:      req.SubmissionID,
		GradedByTeacherID: &teacherID,
		Score:             req.Score,
		Feedback:          req.Feedback,
	}

	if err := s.repo.Create(ctx, grade); err != nil {
		return nil, errorx.Wrap(err, "create grade")
	}

	detailed, err := s.repo.GetByID(ctx, grade.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created grade")
	}

	return s.toResponse(detailed), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*GradeResponse, error) {
	grade, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("grade")
		}
		return nil, errorx.Wrap(err, "get grade")
	}
	return s.toResponse(grade), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string, req *UpdateGradeRequest) error {
	grade, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("grade")
		}
		return errorx.Wrap(err, "get grade")
	}

	// Only admin or the grading teacher can update
	if role != "admin" && (grade.GradedByTeacherID == nil || *grade.GradedByTeacherID != teacherID) {
		return errorx.NewForbiddenError("you can only update grades you created")
	}

	if req.Score != nil {
		grade.Score = *req.Score
	}
	if req.Feedback != nil {
		grade.Feedback = req.Feedback
	}
	grade.GradedByTeacherID = &teacherID

	if err := s.repo.Update(ctx, &grade.Grade); err != nil {
		return errorx.Wrap(err, "update grade")
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string) error {
	grade, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("grade")
		}
		return errorx.Wrap(err, "get grade")
	}

	// Only admin can delete any, teacher can delete own
	if role != "admin" && (grade.GradedByTeacherID == nil || *grade.GradedByTeacherID != teacherID) {
		return errorx.NewForbiddenError("you can only delete grades you created")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errorx.Wrap(err, "delete grade")
	}

	return nil
}

func (s *service) ListByStudent(ctx context.Context, studentID uuid.UUID, page, limit int) (*GradeListResponse, error) {
	offset := (page - 1) * limit
	grades, total, err := s.repo.ListByStudent(ctx, studentID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list student grades")
	}

	return s.toListResponse(grades, total, page, limit), nil
}

func (s *service) ListByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) (*GradeListResponse, error) {
	offset := (page - 1) * limit
	grades, total, err := s.repo.ListByCourse(ctx, courseID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list course grades")
	}

	return s.toListResponse(grades, total, page, limit), nil
}

func (s *service) ListByAssignment(ctx context.Context, assignmentID uuid.UUID, page, limit int) (*GradeListResponse, error) {
	offset := (page - 1) * limit
	grades, total, err := s.repo.ListByAssignment(ctx, assignmentID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list assignment grades")
	}

	return s.toListResponse(grades, total, page, limit), nil
}

func (s *service) GetStudentCourseSummary(ctx context.Context, studentID, courseID uuid.UUID) (*StudentGradeSummary, error) {
	summary, err := s.repo.GetStudentCourseSummary(ctx, studentID, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("course")
		}
		return nil, errorx.Wrap(err, "get student course summary")
	}
	return summary, nil
}

// Helper methods

func (s *service) toResponse(g *GradeWithDetails) *GradeResponse {
	resp := &GradeResponse{
		ID:               g.ID.String(),
		SubmissionID:     g.SubmissionID.String(),
		StudentID:        g.StudentID.String(),
		AssignmentID:     g.AssignmentID.String(),
		Score:            g.Score,
		MaxPoints:        g.MaxPoints,
		Feedback:         g.Feedback,
		AssignmentTitle:  g.AssignmentTitle,
		CourseTitle:      g.CourseTitle,
		StudentFirstName: g.StudentFirstName,
		StudentLastName:  g.StudentLastName,
		StudentEmail:     g.StudentEmail,
		TeacherFirstName: g.TeacherFirstName,
		TeacherLastName:  g.TeacherLastName,
		GradedAt:         g.GradedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if g.GradedByTeacherID != nil {
		id := g.GradedByTeacherID.String()
		resp.GradedByTeacher = &id
	}
	if g.MaxPoints > 0 {
		resp.Percentage = (g.Score / g.MaxPoints) * 100
	}
	return resp
}

func (s *service) toListResponse(grades []GradeWithDetails, total int64, page, limit int) *GradeListResponse {
	var responses []GradeResponse
	for _, g := range grades {
		responses = append(responses, *s.toResponse(&g))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &GradeListResponse{
		Grades:     responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

var _ Service = (*service)(nil)
