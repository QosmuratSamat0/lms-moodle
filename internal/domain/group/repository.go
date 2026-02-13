package group

import "context"

type Repository interface {
	Create(ctx context.Context, group *Group) error
	GetByID(ctx context.Context, id string) (*Group, error)
	GetByCourseID(ctx context.Context, courseID string) ([]*Group, error)
	Update(ctx context.Context, group *Group) error
	Delete(ctx context.Context, id string) error

	AddMember(ctx context.Context, member *GroupMember) error
	RemoveMember(ctx context.Context, groupID, studentID string) error
	GetMembers(ctx context.Context, groupID string) ([]*GroupMember, error)
	GetStudentGroups(ctx context.Context, studentID string) ([]*Group, error)
	IsMember(ctx context.Context, groupID, studentID string) (bool, error)
	CountMembers(ctx context.Context, groupID string) (int, error)
}
