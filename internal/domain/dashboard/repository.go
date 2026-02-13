package dashboard

import "context"

type Repository interface {
	GetStudentDashboard(ctx context.Context, studentID string) (*StudentDashboard, error)
	GetTeacherDashboard(ctx context.Context, teacherID string) (*TeacherDashboard, error)
	GetUpcomingDeadlines(ctx context.Context, studentID string, limit int) ([]UpcomingDeadline, error)
	GetRecentGrades(ctx context.Context, studentID string, limit int) ([]RecentGrade, error)
	GetGradeTrends(ctx context.Context, studentID string, months int) ([]GradeTrend, error)
	GetCourseProgress(ctx context.Context, studentID string) ([]CourseProgress, error)
	CalculateGPA(ctx context.Context, studentID string) (float64, error)
	GetAttendanceRate(ctx context.Context, studentID string) (float64, error)
}
