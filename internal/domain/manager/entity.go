package manager

import "time"

type Manager struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	EmployeeID        string    `json:"employee_id" db:"employee_id"`
	Department        string    `json:"department" db:"department"`
	ManagesCategories []string  `json:"manages_categories" db:"manages_categories"` // JSON array of category_ids
	ManagesTeachers   []string  `json:"manages_teachers" db:"manages_teachers"`     // JSON array of teacher_ids
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`

	// Populated via joins
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName  string `json:"last_name,omitempty" db:"last_name"`
	Email     string `json:"email,omitempty" db:"email"`
}

type ManagerWithDetails struct {
	Manager
	Categories []CategorySummary `json:"categories,omitempty"`
	Teachers   []TeacherSummary  `json:"teachers,omitempty"`
}

type CategorySummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TeacherSummary struct {
	ID         string `json:"id"`
	EmployeeID string `json:"employee_id"`
	FullName   string `json:"full_name"`
}

type CreateManagerInput struct {
	UserID            string   `json:"user_id" binding:"required"`
	EmployeeID        string   `json:"employee_id" binding:"required"`
	Department        string   `json:"department" binding:"required"`
	ManagesCategories []string `json:"manages_categories"`
	ManagesTeachers   []string `json:"manages_teachers"`
}

type UpdateManagerInput struct {
	Department        *string   `json:"department"`
	ManagesCategories *[]string `json:"manages_categories"`
	ManagesTeachers   *[]string `json:"manages_teachers"`
	IsActive          *bool     `json:"is_active"`
}

type ManagerFilter struct {
	Department string `form:"department"`
	IsActive   *bool  `form:"is_active"`
	Limit      int    `form:"limit,default=20"`
	Offset     int    `form:"offset,default=0"`
}
