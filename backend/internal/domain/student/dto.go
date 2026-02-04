package student

import "github.com/google/uuid"

// UpdateStudentRequest represents a request to update student profile
type UpdateStudentRequest struct {
	FirstName *string `json:"first_name" validate:"omitempty,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,max=100"`
	GroupName *string `json:"group_name" validate:"omitempty,max=100"`
}

// StudentResponse represents a student in API responses
type StudentResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	GroupName *string   `json:"group_name"`
	IsActive  bool      `json:"is_active"`
}

// StudentListResponse represents a paginated list of students
type StudentListResponse struct {
	Students   []StudentResponse `json:"students"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

// StudentFilter represents filters for listing students
type StudentFilter struct {
	GroupName *string `form:"group_name"`
	IsActive  *bool   `form:"is_active"`
	Search    *string `form:"search"`
}
