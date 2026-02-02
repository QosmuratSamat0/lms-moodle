package submission

import (
	"time"

	"github.com/google/uuid"
)

// SubmissionStatus represents the status of a submission
type SubmissionStatus string

const (
	StatusDraft       SubmissionStatus = "draft"
	StatusSubmitted   SubmissionStatus = "submitted"
	StatusLate        SubmissionStatus = "late"
	StatusResubmitted SubmissionStatus = "resubmitted"
)

// Submission represents a student submission
type Submission struct {
	ID           uuid.UUID        `json:"id" db:"id"`
	AssignmentID uuid.UUID        `json:"assignment_id" db:"assignment_id"`
	StudentID    uuid.UUID        `json:"student_id" db:"student_id"`
	ContentText  *string          `json:"content_text" db:"content_text"`
	FileURL      *string          `json:"file_url" db:"file_url"`
	SubmittedAt  time.Time        `json:"submitted_at" db:"submitted_at"`
	Status       SubmissionStatus `json:"status" db:"status"`
}

// SubmissionWithDetails includes assignment and student info
type SubmissionWithDetails struct {
	Submission
	AssignmentTitle  string     `json:"assignment_title" db:"assignment_title"`
	CourseTitle      string     `json:"course_title" db:"course_title"`
	StudentFirstName *string    `json:"student_first_name" db:"student_first_name"`
	StudentLastName  *string    `json:"student_last_name" db:"student_last_name"`
	StudentEmail     string     `json:"student_email" db:"student_email"`
	MaxPoints        int        `json:"max_points" db:"max_points"`
	DueAt            *time.Time `json:"due_at" db:"due_at"`
	AllowLate        bool       `json:"allow_late" db:"allow_late"`
	// Grade info if graded
	GradeScore    *int    `json:"grade_score" db:"grade_score"`
	GradeFeedback *string `json:"grade_feedback" db:"grade_feedback"`
}
