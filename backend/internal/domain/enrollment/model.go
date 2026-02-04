package enrollment

import (
	"time"

	"github.com/google/uuid"
)

// EnrollmentStatus represents the status of an enrollment
type EnrollmentStatus string

const (
	StatusActive   EnrollmentStatus = "active"
	StatusPending  EnrollmentStatus = "pending"
	StatusRejected EnrollmentStatus = "rejected"
	StatusDropped  EnrollmentStatus = "dropped"
)

// Enrollment represents a student enrollment in a course
type Enrollment struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	CourseID   uuid.UUID        `json:"course_id" db:"course_id"`
	StudentID  uuid.UUID        `json:"student_id" db:"student_id"`
	Status     EnrollmentStatus `json:"status" db:"status"`
	EnrolledAt time.Time        `json:"enrolled_at" db:"enrolled_at"`
}

// EnrollmentWithDetails includes student and course information
type EnrollmentWithDetails struct {
	Enrollment
	CourseTitle      string  `json:"course_title" db:"course_title"`
	StudentFirstName *string `json:"student_first_name" db:"student_first_name"`
	StudentLastName  *string `json:"student_last_name" db:"student_last_name"`
	StudentEmail     string  `json:"student_email" db:"student_email"`
}
