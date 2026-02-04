package attendance

import "time"

type Attendance struct {
	ID        string    `db:"id"`
	CourseID  string    `db:"course_id"`
	StudentID string    `db:"student_id"`
	Date      time.Time `db:"date"`
	Present   bool      `db:"present"`
	CreatedAt time.Time `db:"created_at"`
}

type CreateAttendanceInput struct {
	CourseID  string
	StudentID string
	Date      time.Time
	Present   bool
}

type Repository interface {
	Create(attendance *Attendance) error
	ListByCourse(courseID string, skip, take int) ([]*Attendance, error)
	ListByStudent(studentID string, skip, take int) ([]*Attendance, error)
	GetByStudentAndDate(studentID string, date time.Time) (*Attendance, error)
	Update(attendance *Attendance) error
	Delete(id string) error
}
