package admin

import "time"

type Admin struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	EmployeeID  string    `json:"employee_id" db:"employee_id"`
	Department  string    `json:"department" db:"department"`
	AccessLevel string    `json:"access_level" db:"access_level"` // "super_admin" | "admin"
	Permissions []string  `json:"permissions" db:"permissions"`   // JSON array
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Populated via joins
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName  string `json:"last_name,omitempty" db:"last_name"`
	Email     string `json:"email,omitempty" db:"email"`
}

type CreateAdminInput struct {
	UserID      string   `json:"user_id" binding:"required"`
	EmployeeID  string   `json:"employee_id" binding:"required"`
	Department  string   `json:"department"`
	AccessLevel string   `json:"access_level" binding:"required"`
	Permissions []string `json:"permissions"`
}

type UpdateAdminInput struct {
	Department  *string   `json:"department"`
	AccessLevel *string   `json:"access_level"`
	Permissions *[]string `json:"permissions"`
	IsActive    *bool     `json:"is_active"`
}

type AdminFilter struct {
	Department  string `form:"department"`
	AccessLevel string `form:"access_level"`
	IsActive    *bool  `form:"is_active"`
	Limit       int    `form:"limit,default=20"`
	Offset      int    `form:"offset,default=0"`
}
