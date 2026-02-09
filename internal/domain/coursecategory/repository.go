package coursecategory

import "context"

type Repository interface {
	Create(ctx context.Context, category *CourseCategory) error
	GetByID(ctx context.Context, id string) (*CourseCategory, error)
	List(ctx context.Context, filter *CourseCategoryFilter) ([]*CourseCategory, int64, error)
	Update(ctx context.Context, category *CourseCategory) error
	Delete(ctx context.Context, id string) error
}
