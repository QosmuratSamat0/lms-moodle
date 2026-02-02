package grade

import "github.com/google/uuid"

// GradeSubmissionRequest represents a request to grade a submission
type GradeSubmissionRequest struct {
	SubmissionID uuid.UUID `json:"submission_id" validate:"required"`
	Score        int       `json:"score" validate:"gte=0"`
	Feedback     *string   `json:"feedback,omitempty" validate:"omitempty,max=5000"`
}

// UpdateGradeRequest represents a request to update a grade
type UpdateGradeRequest struct {
	Score    *int    `json:"score,omitempty" validate:"omitempty,gte=0"`
	Feedback *string `json:"feedback,omitempty" validate:"omitempty,max=5000"`
}

// GradeResponse represents a grade in API responses
type GradeResponse struct {
	ID               string  `json:"id"`
	SubmissionID     string  `json:"submission_id"`
	StudentID        string  `json:"student_id,omitempty"`
	AssignmentID     string  `json:"assignment_id,omitempty"`
	Score            int     `json:"score"`
	MaxPoints        int     `json:"max_points,omitempty"`
	Percentage       float64 `json:"percentage,omitempty"`
	Feedback         *string `json:"feedback,omitempty"`
	GradedByTeacher  *string `json:"graded_by_teacher_id,omitempty"`
	TeacherFirstName *string `json:"teacher_first_name,omitempty"`
	TeacherLastName  *string `json:"teacher_last_name,omitempty"`
	AssignmentTitle  string  `json:"assignment_title,omitempty"`
	CourseTitle      string  `json:"course_title,omitempty"`
	StudentFirstName *string `json:"student_first_name,omitempty"`
	StudentLastName  *string `json:"student_last_name,omitempty"`
	StudentEmail     string  `json:"student_email,omitempty"`
	GradedAt         string  `json:"graded_at"`
}

// GradeListResponse represents a paginated list of grades
type GradeListResponse struct {
	Grades     []GradeResponse `json:"grades"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// StudentGradeSummary represents a student's overall grade summary
type StudentGradeSummary struct {
	StudentID         string  `json:"student_id"`
	CourseID          string  `json:"course_id"`
	CourseTitle       string  `json:"course_title"`
	TotalAssignments  int     `json:"total_assignments"`
	GradedAssignments int     `json:"graded_assignments"`
	TotalPoints       int     `json:"total_points"`
	EarnedPoints      int     `json:"earned_points"`
	AveragePercentage float64 `json:"average_percentage"`
}
