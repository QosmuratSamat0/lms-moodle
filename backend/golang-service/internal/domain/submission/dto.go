package submission

import "github.com/google/uuid"

// CreateSubmissionRequest represents a request to create/submit an assignment
type CreateSubmissionRequest struct {
	AssignmentID uuid.UUID `json:"assignment_id" validate:"required"`
	ContentText  *string   `json:"content_text,omitempty" validate:"omitempty,max=50000"`
	FileURL      *string   `json:"file_url,omitempty" validate:"omitempty,url,max=500"`
}

// UpdateSubmissionRequest represents a request to update a submission
type UpdateSubmissionRequest struct {
	ContentText *string `json:"content_text,omitempty" validate:"omitempty,max=50000"`
	FileURL     *string `json:"file_url,omitempty" validate:"omitempty,url,max=500"`
}

// SubmissionResponse represents a submission in API responses
type SubmissionResponse struct {
	ID               string  `json:"id"`
	AssignmentID     string  `json:"assignment_id"`
	StudentID        string  `json:"student_id"`
	ContentText      *string `json:"content_text,omitempty"`
	FileURL          *string `json:"file_url,omitempty"`
	Status           string  `json:"status"`
	SubmittedAt      string  `json:"submitted_at"`
	AssignmentTitle  string  `json:"assignment_title,omitempty"`
	CourseTitle      string  `json:"course_title,omitempty"`
	StudentFirstName *string `json:"student_first_name,omitempty"`
	StudentLastName  *string `json:"student_last_name,omitempty"`
	StudentEmail     string  `json:"student_email,omitempty"`
	MaxPoints        int     `json:"max_points,omitempty"`
	DueAt            *string `json:"due_at,omitempty"`
	GradeScore       *int    `json:"grade_score,omitempty"`
	GradeFeedback    *string `json:"grade_feedback,omitempty"`
}

// SubmissionListResponse represents a paginated list of submissions
type SubmissionListResponse struct {
	Submissions []SubmissionResponse `json:"submissions"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
	TotalPages  int                  `json:"total_pages"`
}
