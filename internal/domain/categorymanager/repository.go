package categorymanager

import "context"

type Repository interface {
	Create(ctx context.Context, cm *CategoryManager) error
	GetByID(ctx context.Context, id string) (*CategoryManager, error)
	GetByUserAndCategory(ctx context.Context, userID, categoryID string) (*CategoryManager, error)
	GetByUserID(ctx context.Context, userID string) ([]*CategoryManager, error)
	GetByCategoryID(ctx context.Context, categoryID string) ([]*CategoryManager, error)
	List(ctx context.Context, filter *CategoryManagerFilter) ([]*CategoryManager, int64, error)
	Update(ctx context.Context, cm *CategoryManager) error
	Delete(ctx context.Context, id string) error
	GetWithDetails(ctx context.Context, id string) (*CategoryManagerWithDetails, error)
}
