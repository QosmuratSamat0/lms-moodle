package manager

import "context"

type Repository interface {
	Create(ctx context.Context, manager *Manager) error
	GetByID(ctx context.Context, id string) (*Manager, error)
	GetByUserID(ctx context.Context, userID string) (*Manager, error)
	GetByEmployeeID(ctx context.Context, employeeID string) (*Manager, error)
	GetByDepartment(ctx context.Context, department string) ([]*Manager, error)
	List(ctx context.Context, filter *ManagerFilter) ([]*Manager, int64, error)
	Update(ctx context.Context, manager *Manager) error
	Delete(ctx context.Context, id string) error
	AddManagedCategory(ctx context.Context, managerID, categoryID string) error
	RemoveManagedCategory(ctx context.Context, managerID, categoryID string) error
	AddManagedTeacher(ctx context.Context, managerID, teacherID string) error
	RemoveManagedTeacher(ctx context.Context, managerID, teacherID string) error
	GetWithDetails(ctx context.Context, id string) (*ManagerWithDetails, error)
}
