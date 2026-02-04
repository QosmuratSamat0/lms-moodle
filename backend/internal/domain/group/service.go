package group

import (
	"context"
	"errors"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the group service interface
type Service interface {
	Create(ctx context.Context, req *CreateGroupRequest) (*GroupResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*GroupResponse, error)
	GetByCode(ctx context.Context, code string) (*GroupResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateGroupRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, limit int) (*GroupListResponse, error)

	// Teacher assignments
	AssignTeacher(ctx context.Context, req *AssignTeacherRequest) (*TeacherAssignmentResponse, error)
	UnassignTeacher(ctx context.Context, id uuid.UUID) error
	ListTeacherAssignments(ctx context.Context, teacherID uuid.UUID, page, limit int) (*TeacherAssignmentListResponse, error)
	ListGroupAssignments(ctx context.Context, groupID uuid.UUID, page, limit int) (*TeacherAssignmentListResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new group service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *CreateGroupRequest) (*GroupResponse, error) {
	// Check if code already exists
	existing, err := s.repo.GetByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, errorx.NewConflictError("group with this code already exists")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.Wrap(err, "check existing group")
	}

	group := &Group{
		Code:            req.Code,
		Name:            req.Name,
		Description:     req.Description,
		YearOfAdmission: req.YearOfAdmission,
	}

	if err := s.repo.Create(ctx, group); err != nil {
		return nil, errorx.Wrap(err, "create group")
	}

	detailed, err := s.repo.GetByID(ctx, group.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created group")
	}

	return s.toResponse(detailed), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*GroupResponse, error) {
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("group")
		}
		return nil, errorx.Wrap(err, "get group")
	}
	return s.toResponse(group), nil
}

func (s *service) GetByCode(ctx context.Context, code string) (*GroupResponse, error) {
	group, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("group")
		}
		return nil, errorx.Wrap(err, "get group")
	}
	return s.toResponse(group), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req *UpdateGroupRequest) error {
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("group")
		}
		return errorx.Wrap(err, "get group")
	}

	// Apply updates
	if req.Code != nil {
		// Check if new code conflicts
		existing, err := s.repo.GetByCode(ctx, *req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return errorx.NewConflictError("group with this code already exists")
		}
		group.Code = *req.Code
	}
	if req.Name != nil {
		group.Name = req.Name
	}
	if req.Description != nil {
		group.Description = req.Description
	}
	if req.YearOfAdmission != nil {
		group.YearOfAdmission = req.YearOfAdmission
	}

	if err := s.repo.Update(ctx, &group.Group); err != nil {
		return errorx.Wrap(err, "update group")
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("group")
		}
		return errorx.Wrap(err, "get group")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errorx.Wrap(err, "delete group")
	}

	return nil
}

func (s *service) List(ctx context.Context, page, limit int) (*GroupListResponse, error) {
	offset := (page - 1) * limit
	groups, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list groups")
	}
	return s.toListResponse(groups, total, page, limit), nil
}

func (s *service) AssignTeacher(ctx context.Context, req *AssignTeacherRequest) (*TeacherAssignmentResponse, error) {
	// Check if assignment already exists
	existing, err := s.repo.GetTeacherAssignment(ctx, req.TeacherID, req.CourseID, req.GroupID)
	if err == nil && existing != nil {
		return nil, errorx.NewConflictError("teacher is already assigned to this course-group")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.Wrap(err, "check existing assignment")
	}

	assignment := &TeacherCourseGroup{
		TeacherID: req.TeacherID,
		CourseID:  req.CourseID,
		GroupID:   req.GroupID,
	}

	if err := s.repo.AssignTeacher(ctx, assignment); err != nil {
		return nil, errorx.Wrap(err, "assign teacher")
	}

	// Get full details
	assignments, _, err := s.repo.ListTeacherAssignments(ctx, req.TeacherID, 1, 0)
	if err != nil || len(assignments) == 0 {
		// Return basic response
		return &TeacherAssignmentResponse{
			ID:         assignment.ID.String(),
			TeacherID:  assignment.TeacherID.String(),
			CourseID:   assignment.CourseID.String(),
			GroupID:    assignment.GroupID.String(),
			AssignedAt: assignment.AssignedAt.Format("2006-01-02T15:04:05Z07:00"),
		}, nil
	}

	return s.toAssignmentResponse(&assignments[0]), nil
}

func (s *service) UnassignTeacher(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.UnassignTeacher(ctx, id); err != nil {
		return errorx.Wrap(err, "unassign teacher")
	}
	return nil
}

func (s *service) ListTeacherAssignments(ctx context.Context, teacherID uuid.UUID, page, limit int) (*TeacherAssignmentListResponse, error) {
	offset := (page - 1) * limit
	assignments, total, err := s.repo.ListTeacherAssignments(ctx, teacherID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list teacher assignments")
	}
	return s.toAssignmentListResponse(assignments, total, page, limit), nil
}

func (s *service) ListGroupAssignments(ctx context.Context, groupID uuid.UUID, page, limit int) (*TeacherAssignmentListResponse, error) {
	offset := (page - 1) * limit
	assignments, total, err := s.repo.ListGroupAssignments(ctx, groupID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list group assignments")
	}
	return s.toAssignmentListResponse(assignments, total, page, limit), nil
}

// Helper methods

func (s *service) toResponse(g *GroupWithStats) *GroupResponse {
	return &GroupResponse{
		ID:              g.ID.String(),
		Code:            g.Code,
		Name:            g.Name,
		Description:     g.Description,
		YearOfAdmission: g.YearOfAdmission,
		StudentCount:    g.StudentCount,
		CreatedAt:       g.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *service) toListResponse(groups []GroupWithStats, total int64, page, limit int) *GroupListResponse {
	responses := make([]GroupResponse, 0, len(groups))
	for _, g := range groups {
		responses = append(responses, *s.toResponse(&g))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &GroupListResponse{
		Groups:     responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

func (s *service) toAssignmentResponse(a *TeacherCourseGroupDetails) *TeacherAssignmentResponse {
	return &TeacherAssignmentResponse{
		ID:               a.ID.String(),
		TeacherID:        a.TeacherID.String(),
		CourseID:         a.CourseID.String(),
		GroupID:          a.GroupID.String(),
		TeacherFirstName: a.TeacherFirstName,
		TeacherLastName:  a.TeacherLastName,
		CourseTitle:      a.CourseTitle,
		GroupCode:        a.GroupCode,
		AssignedAt:       a.AssignedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *service) toAssignmentListResponse(assignments []TeacherCourseGroupDetails, total int64, page, limit int) *TeacherAssignmentListResponse {
	responses := make([]TeacherAssignmentResponse, 0, len(assignments))
	for _, a := range assignments {
		responses = append(responses, *s.toAssignmentResponse(&a))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &TeacherAssignmentListResponse{
		Assignments: responses,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
	}
}

var _ Service = (*service)(nil)
