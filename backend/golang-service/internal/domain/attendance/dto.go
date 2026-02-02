package attendance

import (
	"github.com/google/uuid"
)

// CreateSessionRequest represents request to create attendance session
type CreateSessionRequest struct {
	CourseID    uuid.UUID `json:"course_id" validate:"required"`
	Title       *string   `json:"title,omitempty" validate:"omitempty,max=200"`
	SessionDate string    `json:"session_date" validate:"required"`
	StartTime   string    `json:"start_time" validate:"required"`
	EndTime     *string   `json:"end_time,omitempty"`
}

// UpdateSessionRequest represents request to update session
type UpdateSessionRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,max=200"`
	SessionDate *string `json:"session_date,omitempty"`
	StartTime   *string `json:"start_time,omitempty"`
	EndTime     *string `json:"end_time,omitempty"`
}

// MarkAttendanceRequest represents request to mark a single student's attendance
type MarkAttendanceRequest struct {
	StudentID uuid.UUID `json:"student_id" validate:"required"`
	Status    string    `json:"status" validate:"required,oneof=present absent late excused"`
	Notes     *string   `json:"notes,omitempty" validate:"omitempty,max=500"`
}

// BulkMarkAttendanceRequest represents request to mark multiple students
type BulkMarkAttendanceRequest struct {
	Marks []MarkAttendanceRequest `json:"marks" validate:"required,dive"`
}

// SessionResponse represents a session in API responses
type SessionResponse struct {
	ID          string  `json:"id"`
	CourseID    string  `json:"course_id"`
	CourseTitle string  `json:"course_title,omitempty"`
	Title       *string `json:"title,omitempty"`
	SessionDate string  `json:"session_date"`
	StartTime   string  `json:"start_time"`
	EndTime     *string `json:"end_time,omitempty"`
	TeacherName string  `json:"teacher_name,omitempty"`
	TotalMarks  int     `json:"total_marks,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// SessionListResponse represents paginated list of sessions
type SessionListResponse struct {
	Sessions   []SessionResponse `json:"sessions"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// MarkResponse represents an attendance mark in API responses
type MarkResponse struct {
	ID               string  `json:"id"`
	SessionID        string  `json:"session_id"`
	StudentID        string  `json:"student_id"`
	StudentFirstName *string `json:"student_first_name,omitempty"`
	StudentLastName  *string `json:"student_last_name,omitempty"`
	StudentEmail     string  `json:"student_email,omitempty"`
	Status           string  `json:"status"`
	Notes            *string `json:"notes,omitempty"`
	MarkedAt         string  `json:"marked_at"`
}

// SessionAttendanceResponse represents a session with all marks
type SessionAttendanceResponse struct {
	Session SessionResponse `json:"session"`
	Marks   []MarkResponse  `json:"marks"`
}

// AttendanceSummaryResponse represents attendance summary
type AttendanceSummaryResponse struct {
	StudentID      string  `json:"student_id"`
	CourseID       string  `json:"course_id"`
	TotalSessions  int     `json:"total_sessions"`
	PresentCount   int     `json:"present_count"`
	AbsentCount    int     `json:"absent_count"`
	LateCount      int     `json:"late_count"`
	ExcusedCount   int     `json:"excused_count"`
	AttendanceRate float64 `json:"attendance_rate"`
}
