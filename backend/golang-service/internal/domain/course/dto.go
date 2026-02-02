package course

import "github.com/google/uuid"

// CreateCourseRequest represents the request to create a course
type CreateCourseRequest struct {
	Title       string  `json:"title" validate:"required,min=3,max=200"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=5000"`
}

// UpdateCourseRequest represents the request to update a course
type UpdateCourseRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=5000"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// TransferCourseRequest represents the request to transfer course ownership
type TransferCourseRequest struct {
	NewOwnerTeacherID uuid.UUID `json:"new_owner_teacher_id" validate:"required"`
}

// CourseResponse represents a course in API responses
type CourseResponse struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Description      *string `json:"description,omitempty"`
	OwnerTeacherID   *string `json:"owner_teacher_id,omitempty"`
	TeacherFirstName *string `json:"teacher_first_name,omitempty"`
	TeacherLastName  *string `json:"teacher_last_name,omitempty"`
	IsActive         bool    `json:"is_active"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// CourseListResponse represents a paginated list of courses
type CourseListResponse struct {
	Courses    []CourseResponse `json:"courses"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// CourseFilter represents filters for course queries
type CourseFilter struct {
	TeacherID *uuid.UUID
	IsActive  *bool
	Search    *string
}
