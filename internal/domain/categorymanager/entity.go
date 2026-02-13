package categorymanager

import "time"

type CategoryManager struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	CategoryID      string    `json:"category_id" db:"category_id"`
	PermissionLevel string    `json:"permission_level" db:"permission_level"` // "view" | "edit" | "admin"
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	// Populated via joins
	FirstName    string `json:"first_name,omitempty" db:"first_name"`
	LastName     string `json:"last_name,omitempty" db:"last_name"`
	Email        string `json:"email,omitempty" db:"email"`
	CategoryName string `json:"category_name,omitempty" db:"category_name"`
}

type CategoryManagerWithDetails struct {
	CategoryManager
	Category     *CategoryInfo `json:"category,omitempty"`
	TotalCourses int           `json:"total_courses"`
}

type CategoryInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryManagerInput struct {
	UserID          string `json:"user_id" binding:"required"`
	CategoryID      string `json:"category_id" binding:"required"`
	PermissionLevel string `json:"permission_level" binding:"required"`
}

type UpdateCategoryManagerInput struct {
	PermissionLevel *string `json:"permission_level"`
	IsActive        *bool   `json:"is_active"`
}

type CategoryManagerFilter struct {
	CategoryID      string `form:"category_id"`
	PermissionLevel string `form:"permission_level"`
	IsActive        *bool  `form:"is_active"`
	Limit           int    `form:"limit,default=20"`
	Offset          int    `form:"offset,default=0"`
}
