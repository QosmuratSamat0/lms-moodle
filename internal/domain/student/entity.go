package student

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/enrollment"
	"github.com/ap1-final-mini-moodle/internal/domain/group"
)

type Student struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	StudentCode string    `json:"student_code" db:"student_code"` // e.g., "STU2024001"
	Major       string    `json:"major" db:"major"`
	Year        int       `json:"year" db:"year"` // 1, 2, 3, 4
	GPA         float64   `json:"gpa" db:"gpa"`
	Status      string    `json:"status" db:"status"` // active, graduated, suspended
	AdmittedAt  time.Time `json:"admitted_at" db:"admitted_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Populated via joins
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName  string `json:"last_name,omitempty" db:"last_name"`
	Email     string `json:"email,omitempty" db:"email"`
}

type StudentWithDetails struct {
	Student
	Enrollments  []*enrollment.Enrollment `json:"enrollments,omitempty"`
	GroupMembers []*group.GroupMember     `json:"groups,omitempty"`
	TotalCourses int                      `json:"total_courses"`
	TotalGroups  int                      `json:"total_groups"`
}

type CreateStudentInput struct {
	UserID      string    `json:"user_id" binding:"required"`
	StudentCode string    `json:"student_code" binding:"required"`
	Major       string    `json:"major" binding:"required"`
	Year        int       `json:"year" binding:"required,min=1,max=6"`
	AdmittedAt  time.Time `json:"admitted_at"`
}

type UpdateStudentInput struct {
	Major  *string  `json:"major"`
	Year   *int     `json:"year"`
	GPA    *float64 `json:"gpa"`
	Status *string  `json:"status"`
}

type StudentFilter struct {
	Major  string `form:"major"`
	Year   int    `form:"year"`
	Status string `form:"status"`
	Limit  int    `form:"limit,default=20"`
	Offset int    `form:"offset,default=0"`
}
