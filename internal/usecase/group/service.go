package group

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/group"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo group.Repository
}

func NewService(repo group.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateGroup(ctx context.Context, input *group.CreateGroupInput) (*group.Group, error) {
	if input.CourseID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.Name == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.MaxStudents <= 0 {
		input.MaxStudents = 30 // default
	}

	g := &group.Group{
		CourseID:    input.CourseID,
		Name:        input.Name,
		Description: input.Description,
		MaxStudents: input.MaxStudents,
	}

	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*group.Group, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	g, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrGroupNotFound
	}

	return g, nil
}

func (s *Service) GetByCourseID(ctx context.Context, courseID string) ([]*group.Group, error) {
	if courseID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetByCourseID(ctx, courseID)
}

func (s *Service) UpdateGroup(ctx context.Context, id string, input *group.UpdateGroupInput) (*group.Group, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	g, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrGroupNotFound
	}

	if input.Name != nil {
		g.Name = *input.Name
	}
	if input.Description != nil {
		g.Description = *input.Description
	}
	if input.MaxStudents != nil {
		if *input.MaxStudents <= 0 {
			return nil, appErrors.ErrInvalidInput
		}
		g.MaxStudents = *input.MaxStudents
	}

	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}

func (s *Service) DeleteGroup(ctx context.Context, id string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrGroupNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) AddMember(ctx context.Context, input *group.AddMemberInput) (*group.GroupMember, error) {
	if input.GroupID == "" || input.StudentID == "" {
		return nil, appErrors.ErrMissingRequired
	}

	// Check if group exists
	g, err := s.repo.GetByID(ctx, input.GroupID)
	if err != nil {
		return nil, appErrors.ErrGroupNotFound
	}

	// Check if already a member
	isMember, _ := s.repo.IsMember(ctx, input.GroupID, input.StudentID)
	if isMember {
		return nil, appErrors.ErrAlreadyInGroup
	}

	// Check if group is full
	count, _ := s.repo.CountMembers(ctx, input.GroupID)
	if count >= g.MaxStudents {
		return nil, appErrors.ErrGroupFull
	}

	member := &group.GroupMember{
		GroupID:   input.GroupID,
		StudentID: input.StudentID,
	}

	if err := s.repo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *Service) RemoveMember(ctx context.Context, groupID, studentID string) error {
	if groupID == "" || studentID == "" {
		return appErrors.ErrMissingRequired
	}

	// Check if member exists
	isMember, _ := s.repo.IsMember(ctx, groupID, studentID)
	if !isMember {
		return appErrors.ErrNotInGroup
	}

	return s.repo.RemoveMember(ctx, groupID, studentID)
}

func (s *Service) GetMembers(ctx context.Context, groupID string) ([]*group.GroupMember, error) {
	if groupID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetMembers(ctx, groupID)
}

func (s *Service) GetStudentGroups(ctx context.Context, studentID string) ([]*group.Group, error) {
	if studentID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetStudentGroups(ctx, studentID)
}
