package course

import "time"

type Course struct {
	ID          string    `db:"id"`
	Code        string    `db:"code"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	TeacherID   string    `db:"teacher_id"`
	MaxPoints   int       `db:"max_points"`
	Active      bool      `db:"active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
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
