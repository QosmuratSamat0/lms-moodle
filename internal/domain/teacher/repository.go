package teacher

import "context"

type Repository interface {
	Create(ctx context.Context, teacher *Teacher) error
	GetByID(ctx context.Context, id string) (*Teacher, error)
	GetByUserID(ctx context.Context, userID string) (*Teacher, error)
	GetByEmployeeID(ctx context.Context, employeeID string) (*Teacher, error)
	List(ctx context.Context, filter *TeacherFilter) ([]*Teacher, int64, error)
	Update(ctx context.Context, teacher *Teacher) error
	Delete(ctx context.Context, id string) error
	GetTeacherCourses(ctx context.Context, teacherID string) ([]TeacherCourse, error)
	GetTeacherGroups(ctx context.Context, teacherID string) ([]*TeacherGroup, error)
}
