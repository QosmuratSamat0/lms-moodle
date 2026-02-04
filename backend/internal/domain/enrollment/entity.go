package enrollment

import "time"

type Enrollment struct {
	ID         string    `db:"id"`
	CourseID   string    `db:"course_id"`
	StudentID  string    `db:"student_id"`
	EnrolledAt time.Time `db:"enrolled_at"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
}

type CreateEnrollmentInput struct {
	CourseID  string
	StudentID string
}

type Repository interface {
	Create(enrollment *Enrollment) error
	GetByID(id string) (*Enrollment, error)
	GetByCourseAndStudent(courseID, studentID string) (*Enrollment, error)
	ListByCourse(courseID string, skip, take int) ([]*Enrollment, error)
	ListByStudent(studentID string, skip, take int) ([]*Enrollment, error)
	Delete(id string) error
}
