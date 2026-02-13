package attendance

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/attendance"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) attendance.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateSession(s *attendance.AttendanceSession) error {
	// Auto-ensure teacher record exists
	if s.CreatedByTeacherID != nil && *s.CreatedByTeacherID != "" {
		_, _ = r.db.Exec(context.Background(),
			`INSERT INTO teachers (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`,
			*s.CreatedByTeacherID)
	}
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO attendance_sessions (id, course_id, starts_at, ends_at, created_by_teacher_id, created_at)
		 VALUES ($1, $2, $3, $4, (SELECT id FROM teachers WHERE user_id = $5), $6)`,
		s.ID, s.CourseID, s.StartsAt, s.EndsAt, s.CreatedByTeacherID, s.CreatedAt)
	return err
}

func (r *PostgresRepository) GetSession(id string) (*attendance.AttendanceSession, error) {
	s := &attendance.AttendanceSession{}
	err := r.db.QueryRow(context.Background(),
		`SELECT s.id, s.course_id, s.starts_at, s.ends_at, s.created_at,
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id), 0),
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id AND m.status = 'present'), 0),
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id AND m.status = 'absent'), 0),
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id AND m.status = 'late'), 0),
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id AND m.status = 'excused'), 0)
		 FROM attendance_sessions s WHERE s.id = $1`, id).
		Scan(&s.ID, &s.CourseID, &s.StartsAt, &s.EndsAt, &s.CreatedAt,
			&s.TotalStudents, &s.PresentCount, &s.AbsentCount, &s.LateCount, &s.ExcusedCount)
	return s, err
}

func (r *PostgresRepository) ListSessionsByCourse(courseID string) ([]*attendance.AttendanceSession, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT s.id, s.course_id, s.starts_at, s.ends_at, s.created_at,
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id), 0),
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id AND m.status = 'present'), 0),
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id AND m.status = 'absent'), 0),
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id AND m.status = 'late'), 0),
		        COALESCE((SELECT COUNT(*) FROM attendance_marks m WHERE m.session_id = s.id AND m.status = 'excused'), 0)
		 FROM attendance_sessions s
		 WHERE s.course_id = $1
		 ORDER BY s.starts_at DESC`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*attendance.AttendanceSession
	for rows.Next() {
		s := &attendance.AttendanceSession{}
		if err := rows.Scan(&s.ID, &s.CourseID, &s.StartsAt, &s.EndsAt, &s.CreatedAt,
			&s.TotalStudents, &s.PresentCount, &s.AbsentCount, &s.LateCount, &s.ExcusedCount); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func (r *PostgresRepository) DeleteSession(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM attendance_sessions WHERE id = $1", id)
	return err
}

func (r *PostgresRepository) UpsertMark(m *attendance.AttendanceMark) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO attendance_marks (id, session_id, student_id, status, marked_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (session_id, student_id) DO UPDATE SET status = $4, marked_at = $5`,
		m.ID, m.SessionID, m.StudentID, m.Status, m.MarkedAt)
	return err
}

func (r *PostgresRepository) BulkUpsertMarks(marks []*attendance.AttendanceMark) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, m := range marks {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO attendance_marks (id, session_id, student_id, status, marked_at)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (session_id, student_id) DO UPDATE SET status = $4, marked_at = $5`,
			m.ID, m.SessionID, m.StudentID, m.Status, m.MarkedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (r *PostgresRepository) GetMarksBySession(sessionID string) ([]*attendance.AttendanceMark, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT m.id, m.session_id, m.student_id, m.status, m.marked_at,
		        COALESCE(st.first_name, u.email), COALESCE(st.last_name, ''), COALESCE(u.email, '')
		 FROM attendance_marks m
		 LEFT JOIN students st ON m.student_id = st.user_id
		 LEFT JOIN users u ON m.student_id = u.id
		 WHERE m.session_id = $1
		 ORDER BY st.last_name, st.first_name`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marks []*attendance.AttendanceMark
	for rows.Next() {
		m := &attendance.AttendanceMark{}
		if err := rows.Scan(&m.ID, &m.SessionID, &m.StudentID, &m.Status, &m.MarkedAt,
			&m.StudentFirstName, &m.StudentLastName, &m.StudentEmail); err != nil {
			return nil, err
		}
		marks = append(marks, m)
	}
	return marks, rows.Err()
}

func (r *PostgresRepository) GetStudentAttendance(studentID, courseID string) ([]*attendance.AttendanceMark, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT m.id, m.session_id, m.student_id, m.status, m.marked_at,
		        '', '', ''
		 FROM attendance_marks m
		 JOIN attendance_sessions s ON m.session_id = s.id
		 WHERE m.student_id = $1 AND s.course_id = $2
		 ORDER BY s.starts_at DESC`, studentID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marks []*attendance.AttendanceMark
	for rows.Next() {
		m := &attendance.AttendanceMark{}
		if err := rows.Scan(&m.ID, &m.SessionID, &m.StudentID, &m.Status, &m.MarkedAt,
			&m.StudentFirstName, &m.StudentLastName, &m.StudentEmail); err != nil {
			return nil, err
		}
		marks = append(marks, m)
	}
	return marks, rows.Err()
}

func (r *PostgresRepository) GetStudentSummary(studentID, courseID string) (*attendance.StudentAttendanceSummary, error) {
	summary := &attendance.StudentAttendanceSummary{CourseID: courseID}
	err := r.db.QueryRow(context.Background(),
		`SELECT
		    COUNT(*) as total,
		    COALESCE(SUM(CASE WHEN m.status = 'present' THEN 1 ELSE 0 END), 0),
		    COALESCE(SUM(CASE WHEN m.status = 'absent' THEN 1 ELSE 0 END), 0),
		    COALESCE(SUM(CASE WHEN m.status = 'late' THEN 1 ELSE 0 END), 0),
		    COALESCE(SUM(CASE WHEN m.status = 'excused' THEN 1 ELSE 0 END), 0)
		 FROM attendance_marks m
		 JOIN attendance_sessions s ON m.session_id = s.id
		 WHERE m.student_id = $1 AND s.course_id = $2`, studentID, courseID).
		Scan(&summary.TotalClasses, &summary.Present, &summary.Absent, &summary.Late, &summary.Excused)
	if err != nil {
		return summary, err
	}
	if summary.TotalClasses > 0 {
		summary.Percentage = float64(summary.Present+summary.Late) / float64(summary.TotalClasses) * 100
	}
	return summary, nil
}

// Unused legacy stubs to satisfy old imports - will be removed
type Attendance = attendance.AttendanceSession
