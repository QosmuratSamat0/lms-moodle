package attendance

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceStatus represents attendance status
type AttendanceStatus string

const (
	StatusPresent AttendanceStatus = "present"
	StatusAbsent  AttendanceStatus = "absent"
	StatusLate    AttendanceStatus = "late"
	StatusExcused AttendanceStatus = "excused"
)

// Session represents an attendance session (class meeting)
type Session struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	CourseID    uuid.UUID  `json:"course_id" db:"course_id"`
	CreatedByID uuid.UUID  `json:"created_by" db:"created_by"`
	Title       *string    `json:"title" db:"title"`
	SessionDate time.Time  `json:"session_date" db:"session_date"`
	StartTime   time.Time  `json:"start_time" db:"start_time"`
	EndTime     *time.Time `json:"end_time" db:"end_time"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// SessionWithDetails includes course info
type SessionWithDetails struct {
	Session
	CourseTitle string `json:"course_title" db:"course_title"`
	TeacherName string `json:"teacher_name" db:"teacher_name"`
	TotalMarks  int    `json:"total_marks" db:"total_marks"`
}

// Mark represents an attendance mark for a student
type Mark struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	SessionID uuid.UUID        `json:"session_id" db:"session_id"`
	StudentID uuid.UUID        `json:"student_id" db:"student_id"`
	Status    AttendanceStatus `json:"status" db:"status"`
	Notes     *string          `json:"notes" db:"notes"`
	MarkedAt  time.Time        `json:"marked_at" db:"marked_at"`
	MarkedBy  uuid.UUID        `json:"marked_by" db:"marked_by"`
}

// MarkWithDetails includes student info
type MarkWithDetails struct {
	Mark
	StudentFirstName *string `json:"student_first_name" db:"student_first_name"`
	StudentLastName  *string `json:"student_last_name" db:"student_last_name"`
	StudentEmail     string  `json:"student_email" db:"student_email"`
}

// StudentAttendanceSummary represents attendance stats for a student in a course
type StudentAttendanceSummary struct {
	StudentID      uuid.UUID `json:"student_id"`
	CourseID       uuid.UUID `json:"course_id"`
	TotalSessions  int       `json:"total_sessions"`
	PresentCount   int       `json:"present_count"`
	AbsentCount    int       `json:"absent_count"`
	LateCount      int       `json:"late_count"`
	ExcusedCount   int       `json:"excused_count"`
	AttendanceRate float64   `json:"attendance_rate"`
}
