package group

import (
	"time"

	"github.com/google/uuid"
)

// Group represents a student group (e.g., SE-2430)
type Group struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Code            string    `json:"code" db:"code"` // e.g., "SE-2430"
	Name            *string   `json:"name" db:"name"`
	Description     *string   `json:"description" db:"description"`
	YearOfAdmission *int      `json:"year_of_admission" db:"year_of_admission"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// GroupWithStats includes student count
type GroupWithStats struct {
	Group
	StudentCount int `json:"student_count" db:"student_count"`
}

// TeacherCourseGroup represents the mapping between teacher, course, and group
type TeacherCourseGroup struct {
	ID         uuid.UUID `json:"id" db:"id"`
	TeacherID  uuid.UUID `json:"teacher_id" db:"teacher_id"`
	CourseID   uuid.UUID `json:"course_id" db:"course_id"`
	GroupID    uuid.UUID `json:"group_id" db:"group_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
}

// TeacherCourseGroupDetails includes related entity names
type TeacherCourseGroupDetails struct {
	TeacherCourseGroup
	TeacherFirstName *string `json:"teacher_first_name" db:"teacher_first_name"`
	TeacherLastName  *string `json:"teacher_last_name" db:"teacher_last_name"`
	CourseTitle      string  `json:"course_title" db:"course_title"`
	GroupCode        string  `json:"group_code" db:"group_code"`
}
