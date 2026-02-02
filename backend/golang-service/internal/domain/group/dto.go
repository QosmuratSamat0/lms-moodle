package group

import "github.com/google/uuid"

// CreateGroupRequest represents a request to create a group
type CreateGroupRequest struct {
	Code            string  `json:"code" validate:"required,min=2,max=50"`
	Name            *string `json:"name,omitempty" validate:"omitempty,max=200"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	YearOfAdmission *int    `json:"year_of_admission,omitempty" validate:"omitempty,min=2000,max=2100"`
}

// UpdateGroupRequest represents a request to update a group
type UpdateGroupRequest struct {
	Code            *string `json:"code,omitempty" validate:"omitempty,min=2,max=50"`
	Name            *string `json:"name,omitempty" validate:"omitempty,max=200"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	YearOfAdmission *int    `json:"year_of_admission,omitempty" validate:"omitempty,min=2000,max=2100"`
}

// AssignTeacherRequest represents a request to assign a teacher to a course-group
type AssignTeacherRequest struct {
	TeacherID uuid.UUID `json:"teacher_id" validate:"required"`
	CourseID  uuid.UUID `json:"course_id" validate:"required"`
	GroupID   uuid.UUID `json:"group_id" validate:"required"`
}

// GroupResponse represents a group in API responses
type GroupResponse struct {
	ID              string  `json:"id"`
	Code            string  `json:"code"`
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	YearOfAdmission *int    `json:"year_of_admission,omitempty"`
	StudentCount    int     `json:"student_count"`
	CreatedAt       string  `json:"created_at"`
}

// GroupListResponse represents a paginated list of groups
type GroupListResponse struct {
	Groups     []GroupResponse `json:"groups"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// TeacherAssignmentResponse represents a teacher-course-group assignment
type TeacherAssignmentResponse struct {
	ID               string  `json:"id"`
	TeacherID        string  `json:"teacher_id"`
	CourseID         string  `json:"course_id"`
	GroupID          string  `json:"group_id"`
	TeacherFirstName *string `json:"teacher_first_name,omitempty"`
	TeacherLastName  *string `json:"teacher_last_name,omitempty"`
	CourseTitle      string  `json:"course_title"`
	GroupCode        string  `json:"group_code"`
	AssignedAt       string  `json:"assigned_at"`
}

// TeacherAssignmentListResponse represents a list of teacher assignments
type TeacherAssignmentListResponse struct {
	Assignments []TeacherAssignmentResponse `json:"assignments"`
	Total       int64                       `json:"total"`
	Page        int                         `json:"page"`
	Limit       int                         `json:"limit"`
	TotalPages  int                         `json:"total_pages"`
}
