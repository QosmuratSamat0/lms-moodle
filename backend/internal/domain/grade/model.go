package grade

import (
	"time"

	"github.com/google/uuid"
)

// Grade represents a grade for a submission
type Grade struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	SubmissionID      uuid.UUID  `json:"submission_id" db:"submission_id"`
	GradedByTeacherID *uuid.UUID `json:"graded_by_teacher_id" db:"graded_by_teacher_id"`
	Score             float64    `json:"score" db:"score"`
	Feedback          *string    `json:"feedback" db:"feedback"`
	GradedAt          time.Time  `json:"graded_at" db:"graded_at"`
}

// GradeWithDetails includes submission and student info
type GradeWithDetails struct {
	Grade
	StudentID        uuid.UUID `json:"student_id" db:"student_id"`
	AssignmentID     uuid.UUID `json:"assignment_id" db:"assignment_id"`
	AssignmentTitle  string    `json:"assignment_title" db:"assignment_title"`
	CourseTitle      string    `json:"course_title" db:"course_title"`
	MaxPoints        float64   `json:"max_points" db:"max_points"`
	StudentFirstName *string   `json:"student_first_name" db:"student_first_name"`
	StudentLastName  *string   `json:"student_last_name" db:"student_last_name"`
	StudentEmail     string    `json:"student_email" db:"student_email"`
	TeacherFirstName *string   `json:"teacher_first_name" db:"teacher_first_name"`
	TeacherLastName  *string   `json:"teacher_last_name" db:"teacher_last_name"`
}
