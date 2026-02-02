package teacher

import (
	"time"

	"github.com/google/uuid"
)

// Teacher represents a teacher profile
type Teacher struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	FirstName  *string   `json:"first_name" db:"first_name"`
	LastName   *string   `json:"last_name" db:"last_name"`
	Department *string   `json:"department" db:"department"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// TeacherWithUser includes user account information
type TeacherWithUser struct {
	Teacher
	Email    string `json:"email" db:"email"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

// TeacherStats represents teacher statistics
type TeacherStats struct {
	UserID           uuid.UUID `json:"user_id"`
	TotalCourses     int       `json:"total_courses"`
	ActiveCourses    int       `json:"active_courses"`
	TotalStudents    int       `json:"total_students"`
	AssignmentsCount int       `json:"assignments_count"`
	GradesGiven      int       `json:"grades_given"`
}

// FullName returns the teacher's full name
func (t *Teacher) FullName() string {
	first, last := "", ""
	if t.FirstName != nil {
		first = *t.FirstName
	}
	if t.LastName != nil {
		last = *t.LastName
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
