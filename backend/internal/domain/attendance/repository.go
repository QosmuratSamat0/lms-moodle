package attendance

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the attendance repository interface
type Repository interface {
	// Session methods
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (*SessionWithDetails, error)
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	ListSessionsByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]SessionWithDetails, int64, error)

	// Mark methods
	CreateMark(ctx context.Context, mark *Mark) error
	UpdateMark(ctx context.Context, mark *Mark) error
	GetMarkBySessionAndStudent(ctx context.Context, sessionID, studentID uuid.UUID) (*Mark, error)
	ListMarksBySession(ctx context.Context, sessionID uuid.UUID) ([]MarkWithDetails, error)
	GetStudentCourseSummary(ctx context.Context, studentID, courseID uuid.UUID) (*StudentAttendanceSummary, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new attendance repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

// Session methods

func (r *repository) CreateSession(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO attendance_sessions (course_id, created_by, title, session_date, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		session.CourseID,
		session.CreatedByID,
		session.Title,
		session.SessionDate,
		session.StartTime,
		session.EndTime,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *repository) GetSessionByID(ctx context.Context, id uuid.UUID) (*SessionWithDetails, error) {
	query := `
		SELECT 
			s.id, s.course_id, s.created_by, s.title, s.session_date, s.start_time, s.end_time, s.created_at,
			c.title AS course_title,
			COALESCE(t.first_name || ' ' || t.last_name, u.email) AS teacher_name,
			(SELECT COUNT(*) FROM attendance_marks WHERE session_id = s.id) AS total_marks
		FROM attendance_sessions s
		JOIN courses c ON s.course_id = c.id
		JOIN users u ON s.created_by = u.id
		LEFT JOIN teachers t ON s.created_by = t.user_id
		WHERE s.id = $1`

	var session SessionWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&session.ID, &session.CourseID, &session.CreatedByID, &session.Title,
		&session.SessionDate, &session.StartTime, &session.EndTime, &session.CreatedAt,
		&session.CourseTitle, &session.TeacherName, &session.TotalMarks,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) UpdateSession(ctx context.Context, session *Session) error {
	query := `
		UPDATE attendance_sessions 
		SET title = $1, session_date = $2, start_time = $3, end_time = $4
		WHERE id = $5`

	_, err := r.db.Exec(ctx, query,
		session.Title, session.SessionDate, session.StartTime, session.EndTime, session.ID,
	)
	return err
}

func (r *repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	// Delete marks first (cascade should handle this, but explicit is clearer)
	_, err := r.db.Exec(ctx, "DELETE FROM attendance_marks WHERE session_id = $1", id)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, "DELETE FROM attendance_sessions WHERE id = $1", id)
	return err
}

func (r *repository) ListSessionsByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]SessionWithDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM attendance_sessions WHERE course_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, courseID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			s.id, s.course_id, s.created_by, s.title, s.session_date, s.start_time, s.end_time, s.created_at,
			c.title AS course_title,
			COALESCE(t.first_name || ' ' || t.last_name, u.email) AS teacher_name,
			(SELECT COUNT(*) FROM attendance_marks WHERE session_id = s.id) AS total_marks
		FROM attendance_sessions s
		JOIN courses c ON s.course_id = c.id
		JOIN users u ON s.created_by = u.id
		LEFT JOIN teachers t ON s.created_by = t.user_id
		WHERE s.course_id = $1
		ORDER BY s.session_date DESC, s.start_time DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanSessions(rows, total)
}

func (r *repository) scanSessions(rows pgx.Rows, total int64) ([]SessionWithDetails, int64, error) {
	var sessions []SessionWithDetails
	for rows.Next() {
		var s SessionWithDetails
		if err := rows.Scan(
			&s.ID, &s.CourseID, &s.CreatedByID, &s.Title,
			&s.SessionDate, &s.StartTime, &s.EndTime, &s.CreatedAt,
			&s.CourseTitle, &s.TeacherName, &s.TotalMarks,
		); err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, s)
	}
	return sessions, total, nil
}

// Mark methods

func (r *repository) CreateMark(ctx context.Context, mark *Mark) error {
	query := `
		INSERT INTO attendance_marks (session_id, student_id, status, notes, marked_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, marked_at`

	return r.db.QueryRow(ctx, query,
		mark.SessionID, mark.StudentID, mark.Status, mark.Notes, mark.MarkedBy,
	).Scan(&mark.ID, &mark.MarkedAt)
}

func (r *repository) UpdateMark(ctx context.Context, mark *Mark) error {
	query := `
		UPDATE attendance_marks 
		SET status = $1, notes = $2, marked_by = $3, marked_at = $4
		WHERE id = $5`

	_, err := r.db.Exec(ctx, query,
		mark.Status, mark.Notes, mark.MarkedBy, time.Now(), mark.ID,
	)
	return err
}

func (r *repository) GetMarkBySessionAndStudent(ctx context.Context, sessionID, studentID uuid.UUID) (*Mark, error) {
	query := `
		SELECT id, session_id, student_id, status, notes, marked_at, marked_by
		FROM attendance_marks
		WHERE session_id = $1 AND student_id = $2`

	var m Mark
	err := r.db.QueryRow(ctx, query, sessionID, studentID).Scan(
		&m.ID, &m.SessionID, &m.StudentID, &m.Status, &m.Notes, &m.MarkedAt, &m.MarkedBy,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) ListMarksBySession(ctx context.Context, sessionID uuid.UUID) ([]MarkWithDetails, error) {
	query := `
		SELECT 
			m.id, m.session_id, m.student_id, m.status, m.notes, m.marked_at, m.marked_by,
			s.first_name AS student_first_name, s.last_name AS student_last_name,
			u.email AS student_email
		FROM attendance_marks m
		JOIN students s ON m.student_id = s.user_id
		JOIN users u ON s.user_id = u.id
		WHERE m.session_id = $1
		ORDER BY s.last_name, s.first_name`

	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marks []MarkWithDetails
	for rows.Next() {
		var m MarkWithDetails
		if err := rows.Scan(
			&m.ID, &m.SessionID, &m.StudentID, &m.Status, &m.Notes, &m.MarkedAt, &m.MarkedBy,
			&m.StudentFirstName, &m.StudentLastName, &m.StudentEmail,
		); err != nil {
			return nil, err
		}
		marks = append(marks, m)
	}
	return marks, nil
}

func (r *repository) GetStudentCourseSummary(ctx context.Context, studentID, courseID uuid.UUID) (*StudentAttendanceSummary, error) {
	query := `
		SELECT 
			COUNT(DISTINCT s.id) AS total_sessions,
			COUNT(CASE WHEN m.status = 'present' THEN 1 END) AS present_count,
			COUNT(CASE WHEN m.status = 'absent' THEN 1 END) AS absent_count,
			COUNT(CASE WHEN m.status = 'late' THEN 1 END) AS late_count,
			COUNT(CASE WHEN m.status = 'excused' THEN 1 END) AS excused_count
		FROM attendance_sessions s
		LEFT JOIN attendance_marks m ON s.id = m.session_id AND m.student_id = $1
		WHERE s.course_id = $2`

	var summary StudentAttendanceSummary
	summary.StudentID = studentID
	summary.CourseID = courseID

	err := r.db.QueryRow(ctx, query, studentID, courseID).Scan(
		&summary.TotalSessions, &summary.PresentCount, &summary.AbsentCount,
		&summary.LateCount, &summary.ExcusedCount,
	)
	if err != nil {
		return nil, err
	}

	// Calculate attendance rate (present + late + excused / total)
	if summary.TotalSessions > 0 {
		attended := summary.PresentCount + summary.LateCount + summary.ExcusedCount
		summary.AttendanceRate = float64(attended) / float64(summary.TotalSessions) * 100
	}

	return &summary, nil
}
