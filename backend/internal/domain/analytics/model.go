package analytics

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents analytics event type
type EventType string

const (
	EventTypePageView     EventType = "page_view"
	EventTypeVideoWatch   EventType = "video_watch"
	EventTypeQuizAttempt  EventType = "quiz_attempt"
	EventTypeAssignment   EventType = "assignment"
	EventTypeLogin        EventType = "login"
	EventTypeLogout       EventType = "logout"
	EventTypeCourseAccess EventType = "course_access"
)

// Event represents an analytics event
type Event struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	CourseID  *uuid.UUID `json:"course_id" db:"course_id"`
	EventType EventType  `json:"event_type" db:"event_type"`
	EventData *string    `json:"event_data" db:"event_data"` // JSON data
	UserAgent *string    `json:"user_agent" db:"user_agent"`
	IPAddress *string    `json:"ip_address" db:"ip_address"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// CourseStat represents course statistics
type CourseStat struct {
	CourseID          uuid.UUID `json:"course_id"`
	CourseTitle       string    `json:"course_title"`
	TotalStudents     int       `json:"total_students"`
	ActiveStudents    int       `json:"active_students"`
	TotalAssignments  int       `json:"total_assignments"`
	TotalSubmissions  int       `json:"total_submissions"`
	AverageGrade      float64   `json:"average_grade"`
	CompletionRate    float64   `json:"completion_rate"`
	AverageAttendance float64   `json:"average_attendance"`
}

// UserStat represents user statistics
type UserStat struct {
	UserID             uuid.UUID  `json:"user_id"`
	TotalLogins        int        `json:"total_logins"`
	LastLogin          *time.Time `json:"last_login"`
	TotalCourseViews   int        `json:"total_course_views"`
	TotalTimeSpent     int        `json:"total_time_spent_minutes"`
	AssignmentsDone    int        `json:"assignments_done"`
	AssignmentsPending int        `json:"assignments_pending"`
}

// SystemStat represents system-wide statistics
type SystemStat struct {
	TotalUsers        int       `json:"total_users"`
	TotalStudents     int       `json:"total_students"`
	TotalTeachers     int       `json:"total_teachers"`
	TotalCourses      int       `json:"total_courses"`
	ActiveCourses     int       `json:"active_courses"`
	TotalEnrollments  int       `json:"total_enrollments"`
	ActiveEnrollments int       `json:"active_enrollments"`
	TotalAssignments  int       `json:"total_assignments"`
	TotalSubmissions  int       `json:"total_submissions"`
	UpdatedAt         time.Time `json:"updated_at"`
}
