package analytics

import "github.com/google/uuid"

// TrackEventRequest represents request to track an event
type TrackEventRequest struct {
	CourseID  *uuid.UUID `json:"course_id,omitempty"`
	EventType string     `json:"event_type" validate:"required,oneof=page_view video_watch quiz_attempt assignment login logout course_access"`
	EventData *string    `json:"event_data,omitempty"`
}

// EventResponse represents an event in API responses
type EventResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	CourseID  *string `json:"course_id,omitempty"`
	EventType string  `json:"event_type"`
	EventData *string `json:"event_data,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// CourseStatResponse represents course statistics response
type CourseStatResponse struct {
	CourseID          string  `json:"course_id"`
	CourseTitle       string  `json:"course_title"`
	TotalStudents     int     `json:"total_students"`
	ActiveStudents    int     `json:"active_students"`
	TotalAssignments  int     `json:"total_assignments"`
	TotalSubmissions  int     `json:"total_submissions"`
	AverageGrade      float64 `json:"average_grade"`
	CompletionRate    float64 `json:"completion_rate"`
	AverageAttendance float64 `json:"average_attendance"`
}

// UserStatResponse represents user statistics response
type UserStatResponse struct {
	UserID             string  `json:"user_id"`
	TotalLogins        int     `json:"total_logins"`
	LastLogin          *string `json:"last_login,omitempty"`
	TotalCourseViews   int     `json:"total_course_views"`
	TotalTimeSpent     int     `json:"total_time_spent_minutes"`
	AssignmentsDone    int     `json:"assignments_done"`
	AssignmentsPending int     `json:"assignments_pending"`
}

// SystemStatResponse represents system statistics response
type SystemStatResponse struct {
	TotalUsers        int    `json:"total_users"`
	TotalStudents     int    `json:"total_students"`
	TotalTeachers     int    `json:"total_teachers"`
	TotalCourses      int    `json:"total_courses"`
	ActiveCourses     int    `json:"active_courses"`
	TotalEnrollments  int    `json:"total_enrollments"`
	ActiveEnrollments int    `json:"active_enrollments"`
	TotalAssignments  int    `json:"total_assignments"`
	TotalSubmissions  int    `json:"total_submissions"`
	UpdatedAt         string `json:"updated_at"`
}

// TimeSeriesPoint represents a data point in time series
type TimeSeriesPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// TimeSeriesResponse represents time series data
type TimeSeriesResponse struct {
	Data []TimeSeriesPoint `json:"data"`
}

// LeaderboardEntry represents a leaderboard entry
type LeaderboardEntry struct {
	Rank         int     `json:"rank"`
	UserID       string  `json:"user_id"`
	UserName     string  `json:"user_name"`
	Score        float64 `json:"score"`
	Assignments  int     `json:"assignments_completed"`
	AverageGrade float64 `json:"average_grade"`
}

// LeaderboardResponse represents leaderboard response
type LeaderboardResponse struct {
	CourseID    string             `json:"course_id"`
	CourseTitle string             `json:"course_title"`
	Entries     []LeaderboardEntry `json:"entries"`
}
