package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the schedule repository interface
type Repository interface {
	Create(ctx context.Context, event *Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*EventWithDetails, error)
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCourse(ctx context.Context, courseID uuid.UUID, filter *EventFilter, limit, offset int) ([]EventWithDetails, int64, error)
	ListByUser(ctx context.Context, userID uuid.UUID, filter *EventFilter, limit, offset int) ([]EventWithDetails, int64, error)
	ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]EventWithDetails, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new schedule repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, event *Event) error {
	query := `
		INSERT INTO schedule_events (
			course_id, created_by, title, description, event_type,
			start_time, end_time, location, is_online, meeting_url,
			recurrence, recur_until
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		event.CourseID, event.CreatedBy, event.Title, event.Description, event.EventType,
		event.StartTime, event.EndTime, event.Location, event.IsOnline, event.MeetingURL,
		event.Recurrence, event.RecurUntil,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*EventWithDetails, error) {
	query := `
		SELECT 
			e.id, e.course_id, e.created_by, e.title, e.description, e.event_type,
			e.start_time, e.end_time, e.location, e.is_online, e.meeting_url,
			e.recurrence, e.recur_until, e.created_at, e.updated_at,
			c.title AS course_title,
			COALESCE(t.first_name || ' ' || t.last_name, u.email) AS creator_name,
			u.email AS creator_email
		FROM schedule_events e
		LEFT JOIN courses c ON e.course_id = c.id
		JOIN users u ON e.created_by = u.id
		LEFT JOIN teachers t ON u.id = t.user_id
		WHERE e.id = $1`

	var event EventWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID, &event.CourseID, &event.CreatedBy, &event.Title, &event.Description, &event.EventType,
		&event.StartTime, &event.EndTime, &event.Location, &event.IsOnline, &event.MeetingURL,
		&event.Recurrence, &event.RecurUntil, &event.CreatedAt, &event.UpdatedAt,
		&event.CourseTitle, &event.CreatorName, &event.CreatorEmail,
	)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *repository) Update(ctx context.Context, event *Event) error {
	query := `
		UPDATE schedule_events SET
			title = $1, description = $2, event_type = $3,
			start_time = $4, end_time = $5, location = $6,
			is_online = $7, meeting_url = $8, recurrence = $9,
			recur_until = $10, updated_at = now()
		WHERE id = $11`

	_, err := r.db.Exec(ctx, query,
		event.Title, event.Description, event.EventType,
		event.StartTime, event.EndTime, event.Location,
		event.IsOnline, event.MeetingURL, event.Recurrence,
		event.RecurUntil, event.ID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM schedule_events WHERE id = $1", id)
	return err
}

func (r *repository) ListByCourse(ctx context.Context, courseID uuid.UUID, filter *EventFilter, limit, offset int) ([]EventWithDetails, int64, error) {
	baseWhere := "WHERE e.course_id = $1"
	args := []interface{}{courseID}
	argIdx := 2

	if filter != nil {
		if filter.EventType != nil {
			baseWhere += fmt.Sprintf(" AND e.event_type = $%d", argIdx)
			args = append(args, *filter.EventType)
			argIdx++
		}
		if filter.StartDate != nil {
			baseWhere += fmt.Sprintf(" AND e.start_time >= $%d", argIdx)
			args = append(args, *filter.StartDate)
			argIdx++
		}
		if filter.EndDate != nil {
			baseWhere += fmt.Sprintf(" AND e.end_time <= $%d", argIdx)
			args = append(args, *filter.EndDate)
			argIdx++
		}
	}

	countQuery := "SELECT COUNT(*) FROM schedule_events e " + baseWhere
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			e.id, e.course_id, e.created_by, e.title, e.description, e.event_type,
			e.start_time, e.end_time, e.location, e.is_online, e.meeting_url,
			e.recurrence, e.recur_until, e.created_at, e.updated_at,
			c.title AS course_title,
			COALESCE(t.first_name || ' ' || t.last_name, u.email) AS creator_name,
			u.email AS creator_email
		FROM schedule_events e
		LEFT JOIN courses c ON e.course_id = c.id
		JOIN users u ON e.created_by = u.id
		LEFT JOIN teachers t ON u.id = t.user_id
		` + baseWhere + fmt.Sprintf(` ORDER BY e.start_time ASC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanEvents(rows, total)
}

func (r *repository) ListByUser(ctx context.Context, userID uuid.UUID, filter *EventFilter, limit, offset int) ([]EventWithDetails, int64, error) {
	// Get events from courses the user is enrolled in or created by the user
	baseWhere := `WHERE (
		e.created_by = $1 
		OR e.course_id IN (SELECT course_id FROM enrollments WHERE student_id = $1 AND status = 'active')
		OR e.course_id IN (SELECT id FROM courses WHERE teacher_id = $1)
	)`
	args := []interface{}{userID}
	argIdx := 2

	if filter != nil {
		if filter.EventType != nil {
			baseWhere += fmt.Sprintf(" AND e.event_type = $%d", argIdx)
			args = append(args, *filter.EventType)
			argIdx++
		}
		if filter.StartDate != nil {
			baseWhere += fmt.Sprintf(" AND e.start_time >= $%d", argIdx)
			args = append(args, *filter.StartDate)
			argIdx++
		}
		if filter.EndDate != nil {
			baseWhere += fmt.Sprintf(" AND e.end_time <= $%d", argIdx)
			args = append(args, *filter.EndDate)
			argIdx++
		}
	}

	countQuery := "SELECT COUNT(*) FROM schedule_events e " + baseWhere
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			e.id, e.course_id, e.created_by, e.title, e.description, e.event_type,
			e.start_time, e.end_time, e.location, e.is_online, e.meeting_url,
			e.recurrence, e.recur_until, e.created_at, e.updated_at,
			c.title AS course_title,
			COALESCE(t.first_name || ' ' || t.last_name, u.email) AS creator_name,
			u.email AS creator_email
		FROM schedule_events e
		LEFT JOIN courses c ON e.course_id = c.id
		JOIN users u ON e.created_by = u.id
		LEFT JOIN teachers t ON u.id = t.user_id
		` + baseWhere + fmt.Sprintf(` ORDER BY e.start_time ASC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanEvents(rows, total)
}

func (r *repository) ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]EventWithDetails, error) {
	query := `
		SELECT 
			e.id, e.course_id, e.created_by, e.title, e.description, e.event_type,
			e.start_time, e.end_time, e.location, e.is_online, e.meeting_url,
			e.recurrence, e.recur_until, e.created_at, e.updated_at,
			c.title AS course_title,
			COALESCE(t.first_name || ' ' || t.last_name, u.email) AS creator_name,
			u.email AS creator_email
		FROM schedule_events e
		LEFT JOIN courses c ON e.course_id = c.id
		JOIN users u ON e.created_by = u.id
		LEFT JOIN teachers t ON u.id = t.user_id
		WHERE (
			e.created_by = $1 
			OR e.course_id IN (SELECT course_id FROM enrollments WHERE student_id = $1 AND status = 'active')
			OR e.course_id IN (SELECT id FROM courses WHERE teacher_id = $1)
		)
		AND e.start_time >= $2 AND e.end_time <= $3
		ORDER BY e.start_time ASC`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events, _, err := r.scanEvents(rows, 0)
	return events, err
}

func (r *repository) scanEvents(rows pgx.Rows, total int64) ([]EventWithDetails, int64, error) {
	var events []EventWithDetails
	for rows.Next() {
		var e EventWithDetails
		if err := rows.Scan(
			&e.ID, &e.CourseID, &e.CreatedBy, &e.Title, &e.Description, &e.EventType,
			&e.StartTime, &e.EndTime, &e.Location, &e.IsOnline, &e.MeetingURL,
			&e.Recurrence, &e.RecurUntil, &e.CreatedAt, &e.UpdatedAt,
			&e.CourseTitle, &e.CreatorName, &e.CreatorEmail,
		); err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}
	return events, total, nil
}
