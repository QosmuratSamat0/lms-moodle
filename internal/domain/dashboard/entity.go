package dashboard

import "time"

type StudentDashboard struct {
	StudentID            string             `json:"student_id"`
	StudentName          string             `json:"student_name"`
	GPA                  float64            `json:"overall_gpa"`
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
	CourseProgress       []CourseProgress   `json:"course_stats"`
}

type UpcomingDeadline struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	CourseID   string    `json:"course_id"`
	CourseName string    `json:"course_title"`
	Title      string    `json:"title"`
	DueDate    time.Time `json:"due_at"`
	DaysLeft   int       `json:"days_remaining"`
	Status     string    `json:"status"`
}

type RecentGrade struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	CourseName string    `json:"course_title"`
	Title      string    `json:"assignment_title"`
	Score      float64   `json:"score"`
	MaxPoints  float64   `json:"max_points"`
	Percentage float64   `json:"percentage"`
	GradedAt   time.Time `json:"graded_at"`
}

type GradeTrend struct {
	Month      string  `json:"month"`
	Year       int     `json:"year"`
	AverageGPA float64 `json:"average_percentage"`
}

type CourseProgress struct {
	CourseID             string  `json:"course_id"`
	CourseName           string  `json:"course_title"`
	CompletedAssignments int     `json:"completed_assignments"`
	TotalAssignments     int     `json:"total_assignments"`
	ProgressPercent      float64 `json:"progress_percentage"`
	CurrentGrade         float64 `json:"current_grade"`
	InstructorName       string  `json:"instructor_name"`
	LetterGrade          string  `json:"letter_grade"`
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
