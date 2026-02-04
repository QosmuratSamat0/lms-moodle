package teacher

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the teacher repository interface
type Repository interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*TeacherWithUser, error)
	List(ctx context.Context, filter *TeacherFilter, page, limit int) ([]TeacherWithUser, int, error)
	Update(ctx context.Context, userID uuid.UUID, req *UpdateTeacherRequest) error
	GetStats(ctx context.Context, userID uuid.UUID) (*TeacherStats, error)
	ListByDepartment(ctx context.Context, department string) ([]TeacherWithUser, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new teacher repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, userID uuid.UUID) (*TeacherWithUser, error) {
	query := `
		SELECT t.user_id, t.first_name, t.last_name, t.department, t.created_at,
		       u.email, u.is_active
		FROM teachers t
		JOIN users u ON u.id = t.user_id
		WHERE t.user_id = $1`

	var teacher TeacherWithUser
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&teacher.UserID, &teacher.FirstName, &teacher.LastName,
		&teacher.Department, &teacher.CreatedAt,
		&teacher.Email, &teacher.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &teacher, nil
}

func (r *repository) List(ctx context.Context, filter *TeacherFilter, page, limit int) ([]TeacherWithUser, int, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter != nil {
		if filter.Department != nil && *filter.Department != "" {
			conditions = append(conditions, fmt.Sprintf("t.department = $%d", argIdx))
			args = append(args, *filter.Department)
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
				"(t.first_name ILIKE $%d OR t.last_name ILIKE $%d OR u.email ILIKE $%d)",
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
		SELECT COUNT(*) FROM teachers t
		JOIN users u ON u.id = t.user_id
		%s`, whereClause)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT t.user_id, t.first_name, t.last_name, t.department, t.created_at,
		       u.email, u.is_active
		FROM teachers t
		JOIN users u ON u.id = t.user_id
		%s
		ORDER BY t.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var teachers []TeacherWithUser
	for rows.Next() {
		var t TeacherWithUser
		if err := rows.Scan(
			&t.UserID, &t.FirstName, &t.LastName, &t.Department, &t.CreatedAt,
			&t.Email, &t.IsActive,
		); err != nil {
			return nil, 0, err
		}
		teachers = append(teachers, t)
	}

	return teachers, total, rows.Err()
}

func (r *repository) Update(ctx context.Context, userID uuid.UUID, req *UpdateTeacherRequest) error {
	query := `
		UPDATE teachers
		SET first_name = COALESCE($2, first_name),
		    last_name = COALESCE($3, last_name),
		    department = COALESCE($4, department)
		WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userID, req.FirstName, req.LastName, req.Department)
	return err
}

func (r *repository) GetStats(ctx context.Context, userID uuid.UUID) (*TeacherStats, error) {
	query := `
		SELECT 
			$1::uuid as user_id,
			(SELECT COUNT(*) FROM courses WHERE owner_teacher_id = $1) as total_courses,
			(SELECT COUNT(*) FROM courses WHERE owner_teacher_id = $1 AND is_active = true) as active_courses,
			(SELECT COUNT(DISTINCT e.student_id) FROM enrollments e
			 JOIN courses c ON c.id = e.course_id
			 WHERE c.owner_teacher_id = $1 AND e.status = 'active') as total_students,
			(SELECT COUNT(*) FROM assignments WHERE created_by_teacher_id = $1) as assignments_count,
			(SELECT COUNT(*) FROM grades WHERE graded_by_teacher_id = $1) as grades_given`

	var stats TeacherStats
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.UserID, &stats.TotalCourses, &stats.ActiveCourses,
		&stats.TotalStudents, &stats.AssignmentsCount, &stats.GradesGiven,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *repository) ListByDepartment(ctx context.Context, department string) ([]TeacherWithUser, error) {
	query := `
		SELECT t.user_id, t.first_name, t.last_name, t.department, t.created_at,
		       u.email, u.is_active
		FROM teachers t
		JOIN users u ON u.id = t.user_id
		WHERE t.department = $1 AND u.is_active = true
		ORDER BY t.last_name, t.first_name`

	rows, err := r.db.Query(ctx, query, department)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []TeacherWithUser
	for rows.Next() {
		var t TeacherWithUser
		if err := rows.Scan(
			&t.UserID, &t.FirstName, &t.LastName, &t.Department, &t.CreatedAt,
			&t.Email, &t.IsActive,
		); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}

	return teachers, rows.Err()
}
