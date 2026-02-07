package admin

import "context"

type Repository interface {
	Create(ctx context.Context, admin *Admin) error
	GetByID(ctx context.Context, id string) (*Admin, error)
	GetByUserID(ctx context.Context, userID string) (*Admin, error)
	GetByEmployeeID(ctx context.Context, employeeID string) (*Admin, error)
	List(ctx context.Context, filter *AdminFilter) ([]*Admin, int64, error)
	Update(ctx context.Context, admin *Admin) error
	Delete(ctx context.Context, id string) error
}
