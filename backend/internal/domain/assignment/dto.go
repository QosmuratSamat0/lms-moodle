package assignment

import (
	"time"

	"github.com/google/uuid"
)

// CreateAssignmentRequest represents a request to create an assignment
type CreateAssignmentRequest struct {
	CourseID    uuid.UUID  `json:"course_id" validate:"required"`
	GroupID     *uuid.UUID `json:"group_id,omitempty"`
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=10000"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	MaxPoints   float64    `json:"max_points" validate:"gte=0,lte=100"`
	AllowLate   bool       `json:"allow_late"`
}

// UpdateAssignmentRequest represents a request to update an assignment
type UpdateAssignmentRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=10000"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	MaxPoints   *float64   `json:"max_points,omitempty" validate:"omitempty,gte=0,lte=100"`
	AllowLate   *bool      `json:"allow_late,omitempty"`
	GroupID     *uuid.UUID `json:"group_id,omitempty"`
}

// AssignmentResponse represents an assignment in API responses
type AssignmentResponse struct {
	ID               string   `json:"id"`
	CourseID         string   `json:"course_id"`
	GroupID          *string  `json:"group_id,omitempty"`
	Title            string   `json:"title"`
	Description      *string  `json:"description,omitempty"`
	DueAt            *string  `json:"due_at,omitempty"`
	MaxPoints        float64  `json:"max_points"`
	AllowLate        bool     `json:"allow_late"`
	CreatedByTeacher *string  `json:"created_by_teacher_id,omitempty"`
	CourseTitle      string   `json:"course_title,omitempty"`
	GroupCode        *string  `json:"group_code,omitempty"`
	TeacherFirstName *string  `json:"teacher_first_name,omitempty"`
	TeacherLastName  *string  `json:"teacher_last_name,omitempty"`
	SubmissionCount  int      `json:"submission_count"`
	CreatedAt        string   `json:"created_at"`
}

// AssignmentListResponse represents a paginated list of assignments
type AssignmentListResponse struct {
	Assignments []AssignmentResponse `json:"assignments"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
	TotalPages  int                  `json:"total_pages"`
}
