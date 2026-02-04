package manager

import (
	"time"

	"github.com/google/uuid"
)

// Manager represents a manager profile
type Manager struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	FirstName *string   `json:"first_name" db:"first_name"`
	LastName  *string   `json:"last_name" db:"last_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ManagerWithUser includes user account information
type ManagerWithUser struct {
	Manager
	Email    string `json:"email" db:"email"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

// SystemOverview represents overall system statistics for managers
type SystemOverview struct {
	TotalUsers      int     `json:"total_users"`
	ActiveUsers     int     `json:"active_users"`
	TotalStudents   int     `json:"total_students"`
	TotalTeachers   int     `json:"total_teachers"`
	TotalManagers   int     `json:"total_managers"`
	TotalCourses    int     `json:"total_courses"`
	ActiveCourses   int     `json:"active_courses"`
	TotalEnrollment int     `json:"total_enrollment"`
	PendingEnroll   int     `json:"pending_enrollments"`
	AvgGrade        float64 `json:"average_grade"`
}

// FullName returns the manager's full name
func (m *Manager) FullName() string {
	first, last := "", ""
	if m.FirstName != nil {
		first = *m.FirstName
	}
	if m.LastName != nil {
		last = *m.LastName
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
