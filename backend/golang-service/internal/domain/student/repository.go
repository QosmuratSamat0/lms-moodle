package student

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the student repository interface
type Repository interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*StudentWithUser, error)
	List(ctx context.Context, filter *StudentFilter, page, limit int) ([]StudentWithUser, int, error)
	Update(ctx context.Context, userID uuid.UUID, req *UpdateStudentRequest) error
	GetStats(ctx context.Context, userID uuid.UUID) (*StudentStats, error)
	ListByGroup(ctx context.Context, groupName string) ([]StudentWithUser, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]StudentWithUser, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new student repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, userID uuid.UUID) (*StudentWithUser, error) {
	query := `
		SELECT s.user_id, s.first_name, s.last_name, s.group_name, s.created_at,
		       u.email, u.is_active
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.user_id = $1`

	var student StudentWithUser
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&student.UserID, &student.FirstName, &student.LastName,
		&student.GroupName, &student.CreatedAt,
		&student.Email, &student.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *repository) List(ctx context.Context, filter *StudentFilter, page, limit int) ([]StudentWithUser, int, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter != nil {
		if filter.GroupName != nil && *filter.GroupName != "" {
			conditions = append(conditions, fmt.Sprintf("s.group_name = $%d", argIdx))
			args = append(args, *filter.GroupName)
			argIdx++
		}
		if filter.IsActive != nil {
			conditions = append(conditions, fmt.Sprintf("u.is_active = $%d", argIdx))
			args = append(args, *filter.IsActive)
			argIdx++
		}
		if filter.Search != nil && *filter.Search != "" {
			searchPattern := "%" + *filter.Search + "%"
			conditions = append(conditions, fmt.Sprintf(
				"(s.first_name ILIKE $%d OR s.last_name ILIKE $%d OR u.email ILIKE $%d)",
				argIdx, argIdx, argIdx,
			))
			args = append(args, searchPattern)
			argIdx++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM students s
		JOIN users u ON u.id = s.user_id
		%s`, whereClause)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT s.user_id, s.first_name, s.last_name, s.group_name, s.created_at,
		       u.email, u.is_active
		FROM students s
		JOIN users u ON u.id = s.user_id
		%s
		ORDER BY s.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []StudentWithUser
	for rows.Next() {
		var s StudentWithUser
		if err := rows.Scan(
			&s.UserID, &s.FirstName, &s.LastName, &s.GroupName, &s.CreatedAt,
			&s.Email, &s.IsActive,
		); err != nil {
			return nil, 0, err
		}
		students = append(students, s)
	}

	return students, total, rows.Err()
}

func (r *repository) Update(ctx context.Context, userID uuid.UUID, req *UpdateStudentRequest) error {
	query := `
		UPDATE students
		SET first_name = COALESCE($2, first_name),
		    last_name = COALESCE($3, last_name),
		    group_name = COALESCE($4, group_name)
		WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userID, req.FirstName, req.LastName, req.GroupName)
	return err
}

func (r *repository) GetStats(ctx context.Context, userID uuid.UUID) (*StudentStats, error) {
	query := `
		SELECT 
			$1::uuid as user_id,
			(SELECT COUNT(*) FROM enrollments WHERE student_id = $1) as total_courses,
			(SELECT COUNT(*) FROM enrollments WHERE student_id = $1 AND status = 'active') as active_courses,
			(SELECT COUNT(*) FROM enrollments WHERE student_id = $1 AND status = 'dropped') as completed_courses,
			COALESCE(
				(SELECT AVG(g.score)::float FROM grades g
				 JOIN submissions s ON s.id = g.submission_id
				 WHERE s.student_id = $1), 0
			) as average_grade,
			(SELECT COUNT(*) FROM submissions WHERE student_id = $1) as submission_count,
			COALESCE(
				(SELECT COUNT(*)::float * 100 / NULLIF((SELECT COUNT(*) FROM attendance_marks WHERE student_id = $1), 0)
				 FROM attendance_marks WHERE student_id = $1 AND status = 'present'), 0
			) as attendance_rate`

	var stats StudentStats
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.UserID, &stats.TotalCourses, &stats.ActiveCourses,
		&stats.CompletedCourses, &stats.AverageGrade,
		&stats.SubmissionCount, &stats.AttendanceRate,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *repository) ListByGroup(ctx context.Context, groupName string) ([]StudentWithUser, error) {
	query := `
		SELECT s.user_id, s.first_name, s.last_name, s.group_name, s.created_at,
		       u.email, u.is_active
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.group_name = $1 AND u.is_active = true
		ORDER BY s.last_name, s.first_name`

	rows, err := r.db.Query(ctx, query, groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []StudentWithUser
	for rows.Next() {
		var s StudentWithUser
		if err := rows.Scan(
			&s.UserID, &s.FirstName, &s.LastName, &s.GroupName, &s.CreatedAt,
			&s.Email, &s.IsActive,
		); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, rows.Err()
}

func (r *repository) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]StudentWithUser, error) {
	query := `
		SELECT s.user_id, s.first_name, s.last_name, s.group_name, s.created_at,
		       u.email, u.is_active
		FROM students s
		JOIN users u ON u.id = s.user_id
		JOIN enrollments e ON e.student_id = s.user_id
		WHERE e.course_id = $1 AND e.status = 'active'
		ORDER BY s.last_name, s.first_name`

	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []StudentWithUser
	for rows.Next() {
		var s StudentWithUser
		if err := rows.Scan(
			&s.UserID, &s.FirstName, &s.LastName, &s.GroupName, &s.CreatedAt,
			&s.Email, &s.IsActive,
		); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, rows.Err()
}
