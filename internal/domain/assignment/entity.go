package assignment

import "time"

type Assignment struct {
	ID                 string     `json:"id" db:"id"`
	CourseID           string     `json:"course_id" db:"course_id"`
	Title              string     `json:"title" db:"title"`
	Description        string     `json:"description" db:"description"`
	MaxPoints          float64    `json:"max_points" db:"max_points"`
	DueAt              *time.Time `json:"due_at" db:"due_at"`
	AllowLate          bool       `json:"allow_late" db:"allow_late"`
	CreatedByTeacherID *string    `json:"created_by_teacher_id" db:"created_by_teacher_id"`
	GradingCategory    string     `json:"grading_category" db:"grading_category"`
	WeightPercentage   float64    `json:"weight_percentage" db:"weight_percentage"`
	FileURL            *string    `json:"file_url" db:"file_url"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	CourseTitle        string     `json:"course_title,omitempty" db:"course_title"`
	TeacherFirstName   string     `json:"teacher_first_name,omitempty" db:"teacher_first_name"`
	TeacherLastName    string     `json:"teacher_last_name,omitempty" db:"teacher_last_name"`
}

type CreateAssignmentInput struct {
	CourseID           string
	Title              string
	Description        string
	MaxPoints          float64
	DueAt              *time.Time
	AllowLate          bool
	CreatedByTeacherID *string
	GradingCategory    string
	WeightPercentage   float64
	FileURL            *string
}

type UpdateAssignmentInput struct {
	Title            *string
	Description      *string
	MaxPoints        *float64
	DueAt            *time.Time
	AllowLate        *bool
	GradingCategory  *string
	WeightPercentage *float64
	FileURL          *string
}

type Repository interface {
	Create(assignment *Assignment) error
	GetByID(id string) (*Assignment, error)
	ListByCourse(courseID string, skip, take int) ([]*Assignment, error)
	Update(assignment *Assignment) error
	Delete(id string) error
}
