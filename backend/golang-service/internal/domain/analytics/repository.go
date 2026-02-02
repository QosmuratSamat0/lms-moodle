package analytics

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the analytics repository interface
type Repository interface {
	TrackEvent(ctx context.Context, event *Event) error
	GetCourseStats(ctx context.Context, courseID uuid.UUID) (*CourseStat, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStat, error)
	GetSystemStats(ctx context.Context) (*SystemStat, error)
	GetEventsByUser(ctx context.Context, userID uuid.UUID, eventType *EventType, limit, offset int) ([]Event, int64, error)
	GetEventsByCourse(ctx context.Context, courseID uuid.UUID, eventType *EventType, limit, offset int) ([]Event, int64, error)
	GetDailyEventCounts(ctx context.Context, eventType EventType, startDate, endDate time.Time) ([]TimeSeriesPoint, error)
	GetCourseLeaderboard(ctx context.Context, courseID uuid.UUID, limit int) ([]LeaderboardEntry, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new analytics repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) TrackEvent(ctx context.Context, event *Event) error {
	query := `
		INSERT INTO analytics_events (user_id, course_id, event_type, event_data, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		event.UserID, event.CourseID, event.EventType, event.EventData, event.UserAgent, event.IPAddress,
	).Scan(&event.ID, &event.CreatedAt)
}

func (r *repository) GetCourseStats(ctx context.Context, courseID uuid.UUID) (*CourseStat, error) {
	query := `
		SELECT
			c.id AS course_id,
			c.title AS course_title,
			(SELECT COUNT(*) FROM enrollments WHERE course_id = c.id) AS total_students,
			(SELECT COUNT(*) FROM enrollments WHERE course_id = c.id AND status = 'active') AS active_students,
			(SELECT COUNT(*) FROM assignments WHERE course_id = c.id) AS total_assignments,
			(SELECT COUNT(*) FROM submissions s JOIN assignments a ON s.assignment_id = a.id WHERE a.course_id = c.id) AS total_submissions,
			COALESCE((SELECT AVG(g.score::float / a.max_points * 100) FROM grades g JOIN submissions s ON g.submission_id = s.id JOIN assignments a ON s.assignment_id = a.id WHERE a.course_id = c.id), 0) AS average_grade,
			COALESCE((SELECT COUNT(DISTINCT s.student_id)::float / NULLIF(COUNT(DISTINCT e.student_id), 0) * 100 FROM submissions s JOIN assignments a ON s.assignment_id = a.id JOIN enrollments e ON e.course_id = a.course_id WHERE a.course_id = c.id), 0) AS completion_rate,
			COALESCE((SELECT COUNT(CASE WHEN am.status IN ('present', 'late') THEN 1 END)::float / NULLIF(COUNT(*), 0) * 100 FROM attendance_marks am JOIN attendance_sessions ass ON am.session_id = ass.id WHERE ass.course_id = c.id), 0) AS average_attendance
		FROM courses c
		WHERE c.id = $1`

	var stat CourseStat
	err := r.db.QueryRow(ctx, query, courseID).Scan(
		&stat.CourseID, &stat.CourseTitle,
		&stat.TotalStudents, &stat.ActiveStudents,
		&stat.TotalAssignments, &stat.TotalSubmissions,
		&stat.AverageGrade, &stat.CompletionRate, &stat.AverageAttendance,
	)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func (r *repository) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStat, error) {
	query := `
		SELECT
			$1::uuid AS user_id,
			(SELECT COUNT(*) FROM analytics_events WHERE user_id = $1 AND event_type = 'login') AS total_logins,
			(SELECT created_at FROM analytics_events WHERE user_id = $1 AND event_type = 'login' ORDER BY created_at DESC LIMIT 1) AS last_login,
			(SELECT COUNT(*) FROM analytics_events WHERE user_id = $1 AND event_type = 'course_access') AS total_course_views,
			0 AS total_time_spent,
			(SELECT COUNT(*) FROM submissions WHERE student_id = $1) AS assignments_done,
			(SELECT COUNT(*) FROM assignments a JOIN enrollments e ON a.course_id = e.course_id WHERE e.student_id = $1 AND e.status = 'active' AND NOT EXISTS (SELECT 1 FROM submissions s WHERE s.assignment_id = a.id AND s.student_id = $1)) AS assignments_pending`

	var stat UserStat
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stat.UserID,
		&stat.TotalLogins, &stat.LastLogin, &stat.TotalCourseViews,
		&stat.TotalTimeSpent, &stat.AssignmentsDone, &stat.AssignmentsPending,
	)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func (r *repository) GetSystemStats(ctx context.Context) (*SystemStat, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM users) AS total_users,
			(SELECT COUNT(*) FROM students) AS total_students,
			(SELECT COUNT(*) FROM teachers) AS total_teachers,
			(SELECT COUNT(*) FROM courses) AS total_courses,
			(SELECT COUNT(*) FROM courses WHERE status = 'active') AS active_courses,
			(SELECT COUNT(*) FROM enrollments) AS total_enrollments,
			(SELECT COUNT(*) FROM enrollments WHERE status = 'active') AS active_enrollments,
			(SELECT COUNT(*) FROM assignments) AS total_assignments,
			(SELECT COUNT(*) FROM submissions) AS total_submissions,
			now() AS updated_at`

	var stat SystemStat
	err := r.db.QueryRow(ctx, query).Scan(
		&stat.TotalUsers, &stat.TotalStudents, &stat.TotalTeachers,
		&stat.TotalCourses, &stat.ActiveCourses,
		&stat.TotalEnrollments, &stat.ActiveEnrollments,
		&stat.TotalAssignments, &stat.TotalSubmissions, &stat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func (r *repository) GetEventsByUser(ctx context.Context, userID uuid.UUID, eventType *EventType, limit, offset int) ([]Event, int64, error) {
	countQuery := `SELECT COUNT(*) FROM analytics_events WHERE user_id = $1`
	listQuery := `
		SELECT id, user_id, course_id, event_type, event_data, user_agent, ip_address, created_at
		FROM analytics_events WHERE user_id = $1`

	args := []interface{}{userID}
	if eventType != nil {
		countQuery += " AND event_type = $2"
		listQuery += " AND event_type = $2"
		args = append(args, *eventType)
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if eventType != nil {
		listQuery += " ORDER BY created_at DESC LIMIT $3 OFFSET $4"
		args = append(args, limit, offset)
	} else {
		listQuery += " ORDER BY created_at DESC LIMIT $2 OFFSET $3"
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.UserID, &e.CourseID, &e.EventType, &e.EventData, &e.UserAgent, &e.IPAddress, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}

	return events, total, nil
}

func (r *repository) GetEventsByCourse(ctx context.Context, courseID uuid.UUID, eventType *EventType, limit, offset int) ([]Event, int64, error) {
	countQuery := `SELECT COUNT(*) FROM analytics_events WHERE course_id = $1`
	listQuery := `
		SELECT id, user_id, course_id, event_type, event_data, user_agent, ip_address, created_at
		FROM analytics_events WHERE course_id = $1`

	args := []interface{}{courseID}
	if eventType != nil {
		countQuery += " AND event_type = $2"
		listQuery += " AND event_type = $2"
		args = append(args, *eventType)
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if eventType != nil {
		listQuery += " ORDER BY created_at DESC LIMIT $3 OFFSET $4"
		args = append(args, limit, offset)
	} else {
		listQuery += " ORDER BY created_at DESC LIMIT $2 OFFSET $3"
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.UserID, &e.CourseID, &e.EventType, &e.EventData, &e.UserAgent, &e.IPAddress, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}

	return events, total, nil
}

func (r *repository) GetDailyEventCounts(ctx context.Context, eventType EventType, startDate, endDate time.Time) ([]TimeSeriesPoint, error) {
	query := `
		SELECT DATE(created_at) AS date, COUNT(*) AS count
		FROM analytics_events
		WHERE event_type = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY DATE(created_at)
		ORDER BY date`

	rows, err := r.db.Query(ctx, query, eventType, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []TimeSeriesPoint
	for rows.Next() {
		var p TimeSeriesPoint
		var date time.Time
		if err := rows.Scan(&date, &p.Count); err != nil {
			return nil, err
		}
		p.Date = date.Format("2006-01-02")
		points = append(points, p)
	}

	return points, nil
}

func (r *repository) GetCourseLeaderboard(ctx context.Context, courseID uuid.UUID, limit int) ([]LeaderboardEntry, error) {
	query := `
		SELECT 
			ROW_NUMBER() OVER (ORDER BY AVG(g.score::float / a.max_points * 100) DESC, COUNT(s.id) DESC) AS rank,
			e.student_id,
			COALESCE(st.first_name || ' ' || st.last_name, u.email) AS user_name,
			COALESCE(AVG(g.score::float / a.max_points * 100), 0) AS score,
			COUNT(DISTINCT s.id) AS assignments_completed,
			COALESCE(AVG(g.score::float / a.max_points * 100), 0) AS average_grade
		FROM enrollments e
		JOIN students st ON e.student_id = st.user_id
		JOIN users u ON st.user_id = u.id
		LEFT JOIN submissions s ON e.student_id = s.student_id
		LEFT JOIN assignments a ON s.assignment_id = a.id AND a.course_id = e.course_id
		LEFT JOIN grades g ON s.id = g.submission_id
		WHERE e.course_id = $1 AND e.status = 'active'
		GROUP BY e.student_id, st.first_name, st.last_name, u.email
		ORDER BY score DESC, assignments_completed DESC
		LIMIT $2`

	rows, err := r.db.Query(ctx, query, courseID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LeaderboardEntry
	for rows.Next() {
		var e LeaderboardEntry
		var userID uuid.UUID
		if err := rows.Scan(&e.Rank, &userID, &e.UserName, &e.Score, &e.Assignments, &e.AverageGrade); err != nil {
			return nil, err
		}
		e.UserID = userID.String()
		entries = append(entries, e)
	}

	return entries, nil
}
