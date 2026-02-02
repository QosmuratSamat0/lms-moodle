package assignment

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the assignment repository interface
type Repository interface {
	Create(ctx context.Context, assignment *Assignment) error
	GetByID(ctx context.Context, id uuid.UUID) (*AssignmentWithDetails, error)
	Update(ctx context.Context, assignment *Assignment) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]AssignmentWithDetails, int64, error)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]AssignmentWithDetails, int64, error)
	ListUpcoming(ctx context.Context, studentID uuid.UUID, limit int) ([]AssignmentWithDetails, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new assignment repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, assignment *Assignment) error {
	query := `
		INSERT INTO assignments (course_id, group_id, title, description, due_at, max_points, allow_late, created_by_teacher_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		assignment.CourseID,
		assignment.GroupID,
		assignment.Title,
		assignment.Description,
		assignment.DueAt,
		assignment.MaxPoints,
		assignment.AllowLate,
		assignment.CreatedByTeacherID,
	).Scan(&assignment.ID, &assignment.CreatedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*AssignmentWithDetails, error) {
	query := `
		SELECT 
			a.id, a.course_id, a.group_id, a.title, a.description, a.due_at, a.max_points, a.allow_late, a.created_by_teacher_id, a.created_at,
			c.title AS course_title,
			g.code AS group_code,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name,
			COALESCE((SELECT COUNT(*) FROM submissions s WHERE s.assignment_id = a.id), 0) AS submission_count
		FROM assignments a
		JOIN courses c ON a.course_id = c.id
		LEFT JOIN groups g ON a.group_id = g.id
		LEFT JOIN teachers t ON a.created_by_teacher_id = t.user_id
		WHERE a.id = $1`

	var a AssignmentWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.CourseID, &a.GroupID, &a.Title, &a.Description, &a.DueAt, &a.MaxPoints, &a.AllowLate, &a.CreatedByTeacherID, &a.CreatedAt,
		&a.CourseTitle,
		&a.GroupCode,
		&a.TeacherFirstName, &a.TeacherLastName,
		&a.SubmissionCount,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) Update(ctx context.Context, assignment *Assignment) error {
	query := `
		UPDATE assignments 
		SET title = $1, description = $2, due_at = $3, max_points = $4, allow_late = $5, group_id = $6
		WHERE id = $7`

	_, err := r.db.Exec(ctx, query,
		assignment.Title, assignment.Description, assignment.DueAt, assignment.MaxPoints, assignment.AllowLate, assignment.GroupID, assignment.ID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM assignments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]AssignmentWithDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM assignments WHERE course_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, courseID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			a.id, a.course_id, a.group_id, a.title, a.description, a.due_at, a.max_points, a.allow_late, a.created_by_teacher_id, a.created_at,
			c.title AS course_title,
			g.code AS group_code,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name,
			COALESCE((SELECT COUNT(*) FROM submissions s WHERE s.assignment_id = a.id), 0) AS submission_count
		FROM assignments a
		JOIN courses c ON a.course_id = c.id
		LEFT JOIN groups g ON a.group_id = g.id
		LEFT JOIN teachers t ON a.created_by_teacher_id = t.user_id
		WHERE a.course_id = $1
		ORDER BY a.due_at ASC NULLS LAST, a.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanAssignments(rows, total)
}

func (r *repository) ListByTeacher(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]AssignmentWithDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM assignments WHERE created_by_teacher_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, teacherID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			a.id, a.course_id, a.group_id, a.title, a.description, a.due_at, a.max_points, a.allow_late, a.created_by_teacher_id, a.created_at,
			c.title AS course_title,
			g.code AS group_code,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name,
			COALESCE((SELECT COUNT(*) FROM submissions s WHERE s.assignment_id = a.id), 0) AS submission_count
		FROM assignments a
		JOIN courses c ON a.course_id = c.id
		LEFT JOIN groups g ON a.group_id = g.id
		LEFT JOIN teachers t ON a.created_by_teacher_id = t.user_id
		WHERE a.created_by_teacher_id = $1
		ORDER BY a.due_at ASC NULLS LAST, a.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, teacherID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanAssignments(rows, total)
}

func (r *repository) ListUpcoming(ctx context.Context, studentID uuid.UUID, limit int) ([]AssignmentWithDetails, error) {
	query := `
		SELECT 
			a.id, a.course_id, a.group_id, a.title, a.description, a.due_at, a.max_points, a.allow_late, a.created_by_teacher_id, a.created_at,
			c.title AS course_title,
			g.code AS group_code,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name,
			0 AS submission_count
		FROM assignments a
		JOIN courses c ON a.course_id = c.id
		JOIN enrollments e ON c.id = e.course_id
		LEFT JOIN groups g ON a.group_id = g.id
		LEFT JOIN teachers t ON a.created_by_teacher_id = t.user_id
		WHERE e.student_id = $1 
			AND e.status = 'active' 
			AND a.due_at > NOW()
			AND NOT EXISTS (SELECT 1 FROM submissions s WHERE s.assignment_id = a.id AND s.student_id = $1)
		ORDER BY a.due_at ASC
		LIMIT $2`

	rows, err := r.db.Query(ctx, query, studentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments, _, err := r.scanAssignments(rows, 0)
	return assignments, err
}

func (r *repository) scanAssignments(rows pgx.Rows, total int64) ([]AssignmentWithDetails, int64, error) {
	var assignments []AssignmentWithDetails
	for rows.Next() {
		var a AssignmentWithDetails
		if err := rows.Scan(
			&a.ID, &a.CourseID, &a.GroupID, &a.Title, &a.Description, &a.DueAt, &a.MaxPoints, &a.AllowLate, &a.CreatedByTeacherID, &a.CreatedAt,
			&a.CourseTitle,
			&a.GroupCode,
			&a.TeacherFirstName, &a.TeacherLastName,
			&a.SubmissionCount,
		); err != nil {
			return nil, 0, err
		}
		assignments = append(assignments, a)
	}
	return assignments, total, nil
}
