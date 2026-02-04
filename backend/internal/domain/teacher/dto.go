package teacher

import "github.com/google/uuid"

// UpdateTeacherRequest represents a request to update teacher profile
type UpdateTeacherRequest struct {
	FirstName  *string `json:"first_name" validate:"omitempty,max=100"`
	LastName   *string `json:"last_name" validate:"omitempty,max=100"`
	Department *string `json:"department" validate:"omitempty,max=100"`
}

// TeacherResponse represents a teacher in API responses
type TeacherResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	FirstName  *string   `json:"first_name"`
	LastName   *string   `json:"last_name"`
	Department *string   `json:"department"`
	IsActive   bool      `json:"is_active"`
}

// TeacherListResponse represents a paginated list of teachers
type TeacherListResponse struct {
	Teachers   []TeacherResponse `json:"teachers"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

// TeacherFilter represents filters for listing teachers
type TeacherFilter struct {
	Department *string `form:"department"`
	IsActive   *bool   `form:"is_active"`
	Search     *string `form:"search"`
}
