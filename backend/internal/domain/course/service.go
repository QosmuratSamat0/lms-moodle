package course

import (
	"context"
	"errors"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the course service interface
type Service interface {
	Create(ctx context.Context, teacherID uuid.UUID, req *CreateCourseRequest) (*CourseResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*CourseResponse, error)
	Update(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string, req *UpdateCourseRequest) error
	Delete(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string) error
	List(ctx context.Context, filter *CourseFilter, page, limit int) (*CourseListResponse, error)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID, page, limit int) (*CourseListResponse, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, page, limit int) (*CourseListResponse, error)
	TransferOwnership(ctx context.Context, courseID, currentOwnerID, newOwnerID uuid.UUID, role string) error
}

type service struct {
	repo Repository
}

// NewService creates a new course service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, teacherID uuid.UUID, req *CreateCourseRequest) (*CourseResponse, error) {
	course := &Course{
		Title:          req.Title,
		Description:    req.Description,
		OwnerTeacherID: &teacherID,
		IsActive:       true,
	}

	if err := s.repo.Create(ctx, course); err != nil {
		return nil, errorx.Wrap(err, "create course")
	}

	// Fetch full course with teacher info
	courseWithTeacher, err := s.repo.GetByID(ctx, course.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created course")
	}

	return s.toCourseResponse(courseWithTeacher), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*CourseResponse, error) {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("course")
		}
		return nil, errorx.Wrap(err, "get course")
	}
	return s.toCourseResponse(course), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string, req *UpdateCourseRequest) error {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("course")
		}
		return errorx.Wrap(err, "get course")
	}

	// Check ownership (admin can update any course)
	if role != "admin" && (course.OwnerTeacherID == nil || *course.OwnerTeacherID != teacherID) {
		return errorx.NewForbiddenError("you don't have permission to update this course")
	}

	// Apply updates
	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.Description != nil {
		course.Description = req.Description
	}
	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, &course.Course); err != nil {
		return errorx.Wrap(err, "update course")
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, teacherID uuid.UUID, role string) error {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("course")
		}
		return errorx.Wrap(err, "get course")
	}

	// Check ownership (admin can delete any course)
	if role != "admin" && (course.OwnerTeacherID == nil || *course.OwnerTeacherID != teacherID) {
		return errorx.NewForbiddenError("you don't have permission to delete this course")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errorx.Wrap(err, "delete course")
	}

	return nil
}

func (s *service) List(ctx context.Context, filter *CourseFilter, page, limit int) (*CourseListResponse, error) {
	offset := (page - 1) * limit
	courses, total, err := s.repo.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list courses")
	}

	return s.toListResponse(courses, total, page, limit), nil
}

func (s *service) ListByTeacher(ctx context.Context, teacherID uuid.UUID, page, limit int) (*CourseListResponse, error) {
	offset := (page - 1) * limit
	courses, total, err := s.repo.ListByTeacher(ctx, teacherID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list teacher courses")
	}

	return s.toListResponse(courses, total, page, limit), nil
}

func (s *service) ListByStudent(ctx context.Context, studentID uuid.UUID, page, limit int) (*CourseListResponse, error) {
	offset := (page - 1) * limit
	courses, total, err := s.repo.ListByStudent(ctx, studentID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list student courses")
	}

	return s.toListResponse(courses, total, page, limit), nil
}

func (s *service) TransferOwnership(ctx context.Context, courseID, currentOwnerID, newOwnerID uuid.UUID, role string) error {
	course, err := s.repo.GetByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("course")
		}
		return errorx.Wrap(err, "get course")
	}

	// Check ownership (admin can transfer any course)
	if role != "admin" && (course.OwnerTeacherID == nil || *course.OwnerTeacherID != currentOwnerID) {
		return errorx.NewForbiddenError("you don't have permission to transfer this course")
	}

	course.OwnerTeacherID = &newOwnerID
	if err := s.repo.Update(ctx, &course.Course); err != nil {
		return errorx.Wrap(err, "transfer course ownership")
	}

	return nil
}

// Helper methods

func (s *service) toCourseResponse(c *CourseWithTeacher) *CourseResponse {
	resp := &CourseResponse{
		ID:               c.ID.String(),
		Title:            c.Title,
		Description:      c.Description,
		IsActive:         c.IsActive,
		TeacherFirstName: c.TeacherFirstName,
		TeacherLastName:  c.TeacherLastName,
		CreatedAt:        c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if c.OwnerTeacherID != nil {
		id := c.OwnerTeacherID.String()
		resp.OwnerTeacherID = &id
	}
	return resp
}

func (s *service) toListResponse(courses []CourseWithTeacher, total int64, page, limit int) *CourseListResponse {
	var responses []CourseResponse
	for _, c := range courses {
		responses = append(responses, *s.toCourseResponse(&c))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &CourseListResponse{
		Courses:    responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

var _ Service = (*service)(nil)
