package attendance

import "time"

// AttendanceSession represents a single attendance-taking event for a course
type AttendanceSession struct {
	ID                 string     `json:"id" db:"id"`
	CourseID           string     `json:"course_id" db:"course_id"`
	StartsAt           time.Time  `json:"starts_at" db:"starts_at"`
	EndsAt             *time.Time `json:"ends_at,omitempty" db:"ends_at"`
	CreatedByTeacherID *string    `json:"created_by_teacher_id,omitempty" db:"created_by_teacher_id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	// Derived
	TotalStudents   int `json:"total_students,omitempty"`
	PresentCount    int `json:"present_count,omitempty"`
	AbsentCount     int `json:"absent_count,omitempty"`
	LateCount       int `json:"late_count,omitempty"`
	ExcusedCount    int `json:"excused_count,omitempty"`
}

// AttendanceMark represents a single student's attendance for a session
type AttendanceMark struct {
	ID               string    `json:"id" db:"id"`
	SessionID        string    `json:"session_id" db:"session_id"`
	StudentID        string    `json:"student_id" db:"student_id"`
	Status           string    `json:"status" db:"status"` // present, absent, late, excused
	MarkedAt         time.Time `json:"marked_at" db:"marked_at"`
	StudentFirstName string    `json:"student_first_name,omitempty"`
	StudentLastName  string    `json:"student_last_name,omitempty"`
	StudentEmail     string    `json:"student_email,omitempty"`
}

// StudentAttendanceSummary is per-course attendance summary for a student
type StudentAttendanceSummary struct {
	CourseID     string  `json:"course_id"`
	TotalClasses int     `json:"total_classes"`
	Present      int     `json:"present"`
	Absent       int     `json:"absent"`
	Late         int     `json:"late"`
	Excused      int     `json:"excused"`
	Percentage   float64 `json:"percentage"`
}

type CreateSessionInput struct {
	CourseID  string
	StartsAt  time.Time
	EndsAt    *time.Time
	TeacherID string
}

type MarkAttendanceInput struct {
	SessionID string
	StudentID string
	Status    string
}

type BulkMarkInput struct {
	SessionID string
	Marks     []SingleMark
}

type SingleMark struct {
	StudentID string `json:"student_id"`
	Status    string `json:"status"`
}

type Repository interface {
	CreateSession(session *AttendanceSession) error
	GetSession(id string) (*AttendanceSession, error)
	ListSessionsByCourse(courseID string) ([]*AttendanceSession, error)
	DeleteSession(id string) error

	UpsertMark(mark *AttendanceMark) error
	BulkUpsertMarks(marks []*AttendanceMark) error
	GetMarksBySession(sessionID string) ([]*AttendanceMark, error)
	GetStudentAttendance(studentID, courseID string) ([]*AttendanceMark, error)
	GetStudentSummary(studentID, courseID string) (*StudentAttendanceSummary, error)
}
