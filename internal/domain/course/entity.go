package course

import "time"

type Course struct {
	ID               string    `db:"id" json:"id"`
	Code             string    `db:"code" json:"code"`
	Title            string    `db:"title" json:"title"`
	Description      string    `db:"description" json:"description"`
	TeacherID        string    `db:"teacher_id" json:"owner_teacher_id"`
	TeacherFirstName string    `db:"-" json:"teacher_first_name"`
	TeacherLastName  string    `db:"-" json:"teacher_last_name"`
	MaxPoints        int       `db:"max_points" json:"max_points"`
	Active           bool      `db:"active" json:"is_active"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type CreateCourseInput struct {
	Code        string
	Title       string
	Description string
	TeacherID   string
	MaxPoints   int
}

type UpdateCourseInput struct {
	Title       *string
	Description *string
	MaxPoints   *int
	Active      *bool
}

type Repository interface {
	Create(course *Course) error
	GetByID(id string) (*Course, error)
	GetByCode(code string) (*Course, error)
	Update(course *Course) error
	List(skip, take int) ([]*Course, error)
	ListByTeacher(teacherID string, skip, take int) ([]*Course, error)
	Delete(id string) error
}
