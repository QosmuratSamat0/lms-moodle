package dashboard

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/dashboard"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo dashboard.Repository
}

func NewService(repo dashboard.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetStudentDashboard(ctx context.Context, studentID string) (*dashboard.StudentDashboard, error) {
	if studentID == "" {
		return nil, appErrors.ErrInvalidID
	}

	// Repository handles complex aggregation
	d, err := s.repo.GetStudentDashboard(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Enrich with additional data
	deadlines, _ := s.repo.GetUpcomingDeadlines(ctx, studentID, 5)
	d.UpcomingDeadlines = deadlines

	recent, _ := s.repo.GetRecentGrades(ctx, studentID, 5)
	d.RecentGrades = recent

	trends, _ := s.repo.GetGradeTrends(ctx, studentID, 6)
	d.GradeTrends = trends

	progress, _ := s.repo.GetCourseProgress(ctx, studentID)
	d.CourseProgress = progress

	return d, nil
}

func (s *Service) GetTeacherDashboard(ctx context.Context, teacherID string) (*dashboard.TeacherDashboard, error) {
	if teacherID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetTeacherDashboard(ctx, teacherID)
}

func (s *Service) GetStudentCourseStats(ctx context.Context, studentID, courseID string) (*dashboard.CourseProgress, error) {
	if studentID == "" || courseID == "" {
		return nil, appErrors.ErrInvalidID
	}

	progressList, err := s.repo.GetCourseProgress(ctx, studentID)
	if err != nil {
		return nil, err
	}

	for i := range progressList {
		if progressList[i].CourseID == courseID {
			return &progressList[i], nil
		}
	}

	return nil, appErrors.ErrEnrollmentNotFound
}
