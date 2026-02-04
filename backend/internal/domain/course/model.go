package course

import (
	"time"

	"github.com/google/uuid"
)

// Course represents a course in the LMS
type Course struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Title          string     `json:"title" db:"title"`
	Description    *string    `json:"description" db:"description"`
	OwnerTeacherID *uuid.UUID `json:"owner_teacher_id" db:"owner_teacher_id"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// CourseWithTeacher includes teacher information
type CourseWithTeacher struct {
	Course
	TeacherFirstName *string `json:"teacher_first_name" db:"teacher_first_name"`
	TeacherLastName  *string `json:"teacher_last_name" db:"teacher_last_name"`
}
