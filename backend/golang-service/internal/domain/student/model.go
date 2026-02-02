package student

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// EnrollmentStatus represents a student's enrollment status
type EnrollmentStatus string

const (
	EnrollmentStatusEnrolled  EnrollmentStatus = "enrolled"
	EnrollmentStatusOnLeave   EnrollmentStatus = "on_leave"
	EnrollmentStatusGraduated EnrollmentStatus = "graduated"
	EnrollmentStatusDropped   EnrollmentStatus = "dropped"
)

// Student represents a student profile
type Student struct {
	UserID           uuid.UUID         `json:"user_id" db:"user_id"`
	FirstName        *string           `json:"first_name" db:"first_name"`
	LastName         *string           `json:"last_name" db:"last_name"`
	GroupID          *uuid.UUID        `json:"group_id" db:"group_id"`
	GroupName        *string           `json:"group_name" db:"group_name"` // Legacy field, keep for compatibility
	Major            *string           `json:"major" db:"major"`
	YearOfStudy      *int              `json:"year_of_study" db:"year_of_study"`
	GPA              *decimal.Decimal  `json:"gpa" db:"gpa"`
	EnrollmentStatus *EnrollmentStatus `json:"enrollment_status" db:"enrollment_status"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
}

// StudentWithUser includes user account information
type StudentWithUser struct {
	Student
	Email     string  `json:"email" db:"email"`
	IsActive  bool    `json:"is_active" db:"is_active"`
	GroupCode *string `json:"group_code" db:"group_code"`
}

// StudentStats represents student statistics
type StudentStats struct {
	UserID           uuid.UUID `json:"user_id"`
	TotalCourses     int       `json:"total_courses"`
	ActiveCourses    int       `json:"active_courses"`
	CompletedCourses int       `json:"completed_courses"`
	AverageGrade     float64   `json:"average_grade"`
	SubmissionCount  int       `json:"submission_count"`
	AttendanceRate   float64   `json:"attendance_rate"`
}

// FullName returns the student's full name
func (s *Student) FullName() string {
	first, last := "", ""
	if s.FirstName != nil {
		first = *s.FirstName
	}
	if s.LastName != nil {
		last = *s.LastName
	}
	if first == "" && last == "" {
		return "Unknown"
	}
	if first == "" {
		return last
	}
	if last == "" {
		return first
	}
	return first + " " + last
}
