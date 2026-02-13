package announcement

import "context"

type Repository interface {
	Create(ctx context.Context, announcement *Announcement) error
	GetByID(ctx context.Context, id string) (*Announcement, error)
	GetByCourseID(ctx context.Context, courseID string, limit, offset int) ([]*Announcement, int64, error)
	Update(ctx context.Context, announcement *Announcement) error
	Delete(ctx context.Context, id string) error
}
