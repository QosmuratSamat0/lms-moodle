package assignment

import "time"

type Assignment struct {
	ID          string    `db:"id"`
	CourseID    string    `db:"course_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	MaxPoints   int       `db:"max_points"`
	DueDate     time.Time `db:"due_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateAssignmentInput struct {
	CourseID    string
	Title       string
	Description string
	MaxPoints   int
	DueDate     time.Time
}

type UpdateAssignmentInput struct {
	Title       *string
	Description *string
	MaxPoints   *int
	DueDate     *time.Time
}

type Repository interface {
	Create(assignment *Assignment) error
	GetByID(id string) (*Assignment, error)
	ListByCourse(courseID string, skip, take int) ([]*Assignment, error)
	Update(assignment *Assignment) error
	Delete(id string) error
}
