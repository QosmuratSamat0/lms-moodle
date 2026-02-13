package enrollment

import "time"

type Enrollment struct {
	ID         string    `db:"id" json:"id"`
	CourseID   string    `db:"course_id" json:"course_id"`
	StudentID  string    `db:"student_id" json:"student_id"`
	EnrolledAt time.Time `db:"enrolled_at" json:"enrolled_at"`
	Status     string    `db:"status" json:"status"`
	// Derived from JOIN
	FirstName  string    `json:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Email      string    `json:"email,omitempty"`
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
