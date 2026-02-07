package appeal

import "context"

type Repository interface {
	Create(ctx context.Context, appeal *GradeAppeal) error
	GetByID(ctx context.Context, id string) (*GradeAppeal, error)
	GetByGradeID(ctx context.Context, gradeID string) (*GradeAppeal, error)
	GetByStudentID(ctx context.Context, studentID string, status *AppealStatus, limit, offset int) ([]*AppealWithDetails, int64, error)
	GetByTeacherCourses(ctx context.Context, teacherID string, status *AppealStatus, limit, offset int) ([]*AppealWithDetails, int64, error)
	Update(ctx context.Context, appeal *GradeAppeal) error
	Delete(ctx context.Context, id string) error
	HasPendingAppeal(ctx context.Context, gradeID string) (bool, error)
}
