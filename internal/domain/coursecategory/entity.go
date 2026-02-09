package coursecategory

import "time"

type CourseCategory struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	Order       int       `json:"order" db:"order"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Populated via joins
	TotalCourses int `json:"total_courses,omitempty"`
}

type CreateCourseCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Order       int    `json:"order"`
}

type UpdateCourseCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Order       *int    `json:"order"`
	IsActive    *bool   `json:"is_active"`
}

type CourseCategoryFilter struct {
	IsActive *bool `form:"is_active"`
	Limit    int   `form:"limit,default=20"`
	Offset   int   `form:"offset,default=0"`
}
