package assignment

import (
	"time"

	"github.com/google/uuid"
)

// Assignment represents a course assignment
type Assignment struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	CourseID           uuid.UUID  `json:"course_id" db:"course_id"`
	GroupID            *uuid.UUID `json:"group_id" db:"group_id"`
	Title              string     `json:"title" db:"title"`
	Description        *string    `json:"description" db:"description"`
	DueAt              *time.Time `json:"due_at" db:"due_at"`
	MaxPoints          float64    `json:"max_points" db:"max_points"`
	AllowLate          bool       `json:"allow_late" db:"allow_late"`
	CreatedByTeacherID *uuid.UUID `json:"created_by_teacher_id" db:"created_by_teacher_id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
}

// AssignmentWithDetails includes course and teacher info
type AssignmentWithDetails struct {
	Assignment
	CourseTitle      string  `json:"course_title" db:"course_title"`
	GroupCode        *string `json:"group_code" db:"group_code"`
	TeacherFirstName *string `json:"teacher_first_name" db:"teacher_first_name"`
	TeacherLastName  *string `json:"teacher_last_name" db:"teacher_last_name"`
	SubmissionCount  int     `json:"submission_count" db:"submission_count"`
}
