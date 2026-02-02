package course

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the course repository interface
type Repository interface {
	Create(ctx context.Context, course *Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*CourseWithTeacher, error)
	Update(ctx context.Context, course *Course) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *CourseFilter, limit, offset int) ([]CourseWithTeacher, int64, error)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]CourseWithTeacher, int64, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]CourseWithTeacher, int64, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new course repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, course *Course) error {
	query := `
		INSERT INTO courses (title, description, owner_teacher_id, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		course.Title,
		course.Description,
		course.OwnerTeacherID,
		course.IsActive,
	).Scan(&course.ID, &course.CreatedAt, &course.UpdatedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*CourseWithTeacher, error) {
	query := `
		SELECT 
			c.id, c.title, c.description, c.owner_teacher_id, c.is_active, c.created_at, c.updated_at,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name
		FROM courses c
		LEFT JOIN teachers t ON c.owner_teacher_id = t.user_id
		WHERE c.id = $1`

	var course CourseWithTeacher
	err := r.db.QueryRow(ctx, query, id).Scan(
		&course.ID, &course.Title, &course.Description,
		&course.OwnerTeacherID, &course.IsActive,
		&course.CreatedAt, &course.UpdatedAt,
		&course.TeacherFirstName, &course.TeacherLastName,
	)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *repository) Update(ctx context.Context, course *Course) error {
	query := `
		UPDATE courses 
		SET title = $1, description = $2, owner_teacher_id = $3, is_active = $4, updated_at = now()
		WHERE id = $5`

	_, err := r.db.Exec(ctx, query,
		course.Title, course.Description, course.OwnerTeacherID, course.IsActive, course.ID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM courses WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) List(ctx context.Context, filter *CourseFilter, limit, offset int) ([]CourseWithTeacher, int64, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter != nil {
		if filter.TeacherID != nil {
			conditions = append(conditions, fmt.Sprintf("c.owner_teacher_id = $%d", argIdx))
			args = append(args, *filter.TeacherID)
			argIdx++
		}
		if filter.IsActive != nil {
			conditions = append(conditions, fmt.Sprintf("c.is_active = $%d", argIdx))
			args = append(args, *filter.IsActive)
			argIdx++
		}
		if filter.Search != nil && *filter.Search != "" {
			conditions = append(conditions, fmt.Sprintf("(c.title ILIKE $%d OR c.description ILIKE $%d)", argIdx, argIdx))
			args = append(args, "%"+*filter.Search+"%")
			argIdx++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM courses c %s`, whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT 
			c.id, c.title, c.description, c.owner_teacher_id, c.is_active, c.created_at, c.updated_at,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name
		FROM courses c
		LEFT JOIN teachers t ON c.owner_teacher_id = t.user_id
		%s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanCourses(rows, total)
}

func (r *repository) ListByTeacher(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]CourseWithTeacher, int64, error) {
	filter := &CourseFilter{TeacherID: &teacherID}
	return r.List(ctx, filter, limit, offset)
}

func (r *repository) ListByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]CourseWithTeacher, int64, error) {
	countQuery := `
		SELECT COUNT(*) 
		FROM courses c
		JOIN enrollments e ON c.id = e.course_id
		WHERE e.student_id = $1 AND e.status = 'active'`

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, studentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			c.id, c.title, c.description, c.owner_teacher_id, c.is_active, c.created_at, c.updated_at,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name
		FROM courses c
		JOIN enrollments e ON c.id = e.course_id
		LEFT JOIN teachers t ON c.owner_teacher_id = t.user_id
		WHERE e.student_id = $1 AND e.status = 'active'
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanCourses(rows, total)
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r *repository) scanCourses(rows pgx.Rows, total int64) ([]CourseWithTeacher, int64, error) {
	var courses []CourseWithTeacher
	for rows.Next() {
		var c CourseWithTeacher
		if err := rows.Scan(
			&c.ID, &c.Title, &c.Description,
			&c.OwnerTeacherID, &c.IsActive,
			&c.CreatedAt, &c.UpdatedAt,
			&c.TeacherFirstName, &c.TeacherLastName,
		); err != nil {
			return nil, 0, err
		}
		courses = append(courses, c)
	}
	return courses, total, nil
}
