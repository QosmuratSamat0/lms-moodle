package analytics

import (
	"context"
	"errors"
	"time"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the analytics service interface
type Service interface {
	TrackEvent(ctx context.Context, userID uuid.UUID, req *TrackEventRequest, userAgent, ipAddress *string) error
	GetCourseStats(ctx context.Context, courseID uuid.UUID) (*CourseStatResponse, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStatResponse, error)
	GetSystemStats(ctx context.Context) (*SystemStatResponse, error)
	GetDailyLogins(ctx context.Context, startDate, endDate time.Time) (*TimeSeriesResponse, error)
	GetCourseLeaderboard(ctx context.Context, courseID uuid.UUID, limit int) (*LeaderboardResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new analytics service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) TrackEvent(ctx context.Context, userID uuid.UUID, req *TrackEventRequest, userAgent, ipAddress *string) error {
	event := &Event{
		UserID:    userID,
		CourseID:  req.CourseID,
		EventType: EventType(req.EventType),
		EventData: req.EventData,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	if err := s.repo.TrackEvent(ctx, event); err != nil {
		return errorx.Wrap(err, "track event")
	}

	return nil
}

func (s *service) GetCourseStats(ctx context.Context, courseID uuid.UUID) (*CourseStatResponse, error) {
	stat, err := s.repo.GetCourseStats(ctx, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("course")
		}
		return nil, errorx.Wrap(err, "get course stats")
	}

	return &CourseStatResponse{
		CourseID:          stat.CourseID.String(),
		CourseTitle:       stat.CourseTitle,
		TotalStudents:     stat.TotalStudents,
		ActiveStudents:    stat.ActiveStudents,
		TotalAssignments:  stat.TotalAssignments,
		TotalSubmissions:  stat.TotalSubmissions,
		AverageGrade:      stat.AverageGrade,
		CompletionRate:    stat.CompletionRate,
		AverageAttendance: stat.AverageAttendance,
	}, nil
}

func (s *service) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStatResponse, error) {
	stat, err := s.repo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, errorx.Wrap(err, "get user stats")
	}

	resp := &UserStatResponse{
		UserID:             stat.UserID.String(),
		TotalLogins:        stat.TotalLogins,
		TotalCourseViews:   stat.TotalCourseViews,
		TotalTimeSpent:     stat.TotalTimeSpent,
		AssignmentsDone:    stat.AssignmentsDone,
		AssignmentsPending: stat.AssignmentsPending,
	}
	if stat.LastLogin != nil {
		ll := stat.LastLogin.Format(time.RFC3339)
		resp.LastLogin = &ll
	}

	return resp, nil
}

func (s *service) GetSystemStats(ctx context.Context) (*SystemStatResponse, error) {
	stat, err := s.repo.GetSystemStats(ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "get system stats")
	}

	return &SystemStatResponse{
		TotalUsers:        stat.TotalUsers,
		TotalStudents:     stat.TotalStudents,
		TotalTeachers:     stat.TotalTeachers,
		TotalCourses:      stat.TotalCourses,
		ActiveCourses:     stat.ActiveCourses,
		TotalEnrollments:  stat.TotalEnrollments,
		ActiveEnrollments: stat.ActiveEnrollments,
		TotalAssignments:  stat.TotalAssignments,
		TotalSubmissions:  stat.TotalSubmissions,
		UpdatedAt:         stat.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) GetDailyLogins(ctx context.Context, startDate, endDate time.Time) (*TimeSeriesResponse, error) {
	points, err := s.repo.GetDailyEventCounts(ctx, EventTypeLogin, startDate, endDate)
	if err != nil {
		return nil, errorx.Wrap(err, "get daily logins")
	}

	return &TimeSeriesResponse{Data: points}, nil
}

func (s *service) GetCourseLeaderboard(ctx context.Context, courseID uuid.UUID, limit int) (*LeaderboardResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Get course info
	stat, err := s.repo.GetCourseStats(ctx, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("course")
		}
		return nil, errorx.Wrap(err, "get course")
	}

	entries, err := s.repo.GetCourseLeaderboard(ctx, courseID, limit)
	if err != nil {
		return nil, errorx.Wrap(err, "get leaderboard")
	}

	return &LeaderboardResponse{
		CourseID:    courseID.String(),
		CourseTitle: stat.CourseTitle,
		Entries:     entries,
	}, nil
}

var _ Service = (*service)(nil)
