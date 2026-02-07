package dashboard

import "time"

type StudentDashboard struct {
	StudentID            string             `json:"student_id"`
	StudentName          string             `json:"student_name"`
	GPA                  float64            `json:"gpa"`
	TotalCourses         int                `json:"total_courses"`
	ActiveCourses        int                `json:"active_courses"`
	CompletedAssignments int                `json:"completed_assignments"`
	PendingAssignments   int                `json:"pending_assignments"`
	TotalQuizzes         int                `json:"total_quizzes"`
	AverageQuizScore     float64            `json:"average_quiz_score"`
	AttendanceRate       float64            `json:"attendance_rate"`
	UpcomingDeadlines    []UpcomingDeadline `json:"upcoming_deadlines"`
	RecentGrades         []RecentGrade      `json:"recent_grades"`
	GradeTrends          []GradeTrend       `json:"grade_trends"`
	CourseProgress       []CourseProgress   `json:"course_progress"`
}

type UpcomingDeadline struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // assignment, quiz
	CourseID   string    `json:"course_id"`
	CourseName string    `json:"course_name"`
	Title      string    `json:"title"`
	DueDate    time.Time `json:"due_date"`
	DaysLeft   int       `json:"days_left"`
}

type RecentGrade struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // assignment, quiz
	CourseName string    `json:"course_name"`
	Title      string    `json:"title"`
	Score      float64   `json:"score"`
	MaxPoints  float64   `json:"max_points"`
	Percentage float64   `json:"percentage"`
	GradedAt   time.Time `json:"graded_at"`
}

type GradeTrend struct {
	Month      string  `json:"month"`
	Year       int     `json:"year"`
	AverageGPA float64 `json:"average_gpa"`
}

type CourseProgress struct {
	CourseID             string  `json:"course_id"`
	CourseName           string  `json:"course_name"`
	CompletedAssignments int     `json:"completed_assignments"`
	TotalAssignments     int     `json:"total_assignments"`
	ProgressPercent      float64 `json:"progress_percent"`
	CurrentGrade         float64 `json:"current_grade"`
}

type TeacherDashboard struct {
	TeacherID          string             `json:"teacher_id"`
	TeacherName        string             `json:"teacher_name"`
	TotalCourses       int                `json:"total_courses"`
	TotalStudents      int                `json:"total_students"`
	PendingSubmissions int                `json:"pending_submissions"`
	PendingAppeals     int                `json:"pending_appeals"`
	RecentSubmissions  []RecentSubmission `json:"recent_submissions"`
	CourseStats        []CourseStats      `json:"course_stats"`
}

type RecentSubmission struct {
	ID              string    `json:"id"`
	StudentName     string    `json:"student_name"`
	CourseID        string    `json:"course_id"`
	CourseName      string    `json:"course_name"`
	AssignmentTitle string    `json:"assignment_title"`
	SubmittedAt     time.Time `json:"submitted_at"`
	Graded          bool      `json:"graded"`
}

type CourseStats struct {
	CourseID       string  `json:"course_id"`
	CourseName     string  `json:"course_name"`
	StudentCount   int     `json:"student_count"`
	AverageGrade   float64 `json:"average_grade"`
	AttendanceRate float64 `json:"attendance_rate"`
	SubmissionRate float64 `json:"submission_rate"`
}
