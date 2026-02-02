package enrollment

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the enrollment repository interface
type Repository interface {
	Create(ctx context.Context, enrollment *Enrollment) error
	GetByID(ctx context.Context, id uuid.UUID) (*EnrollmentWithDetails, error)
	GetByCourseAndStudent(ctx context.Context, courseID, studentID uuid.UUID) (*Enrollment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status EnrollmentStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]EnrollmentWithDetails, int64, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]EnrollmentWithDetails, int64, error)
	IsEnrolled(ctx context.Context, courseID, studentID uuid.UUID) (bool, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new enrollment repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, enrollment *Enrollment) error {
	query := `
		INSERT INTO enrollments (course_id, student_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, enrolled_at`

	return r.db.QueryRow(ctx, query,
		enrollment.CourseID,
		enrollment.StudentID,
		enrollment.Status,
	).Scan(&enrollment.ID, &enrollment.EnrolledAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*EnrollmentWithDetails, error) {
	query := `
		SELECT 
			e.id, e.course_id, e.student_id, e.status, e.enrolled_at,
			c.title AS course_title,
			s.first_name AS student_first_name, s.last_name AS student_last_name,
			u.email AS student_email
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		JOIN students s ON e.student_id = s.user_id
		JOIN users u ON s.user_id = u.id
		WHERE e.id = $1`

	var e EnrollmentWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.CourseID, &e.StudentID, &e.Status, &e.EnrolledAt,
		&e.CourseTitle,
		&e.StudentFirstName, &e.StudentLastName, &e.StudentEmail,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) GetByCourseAndStudent(ctx context.Context, courseID, studentID uuid.UUID) (*Enrollment, error) {
	query := `
		SELECT id, course_id, student_id, status, enrolled_at
		FROM enrollments
		WHERE course_id = $1 AND student_id = $2`

	var e Enrollment
	err := r.db.QueryRow(ctx, query, courseID, studentID).Scan(
		&e.ID, &e.CourseID, &e.StudentID, &e.Status, &e.EnrolledAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status EnrollmentStatus) error {
	query := `UPDATE enrollments SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM enrollments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]EnrollmentWithDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM enrollments WHERE course_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, courseID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			e.id, e.course_id, e.student_id, e.status, e.enrolled_at,
			c.title AS course_title,
			s.first_name AS student_first_name, s.last_name AS student_last_name,
			u.email AS student_email
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		JOIN students s ON e.student_id = s.user_id
		JOIN users u ON s.user_id = u.id
		WHERE e.course_id = $1
		ORDER BY e.enrolled_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanEnrollments(rows, total)
}

func (r *repository) ListByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]EnrollmentWithDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM enrollments WHERE student_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, studentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			e.id, e.course_id, e.student_id, e.status, e.enrolled_at,
			c.title AS course_title,
			s.first_name AS student_first_name, s.last_name AS student_last_name,
			u.email AS student_email
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		JOIN students s ON e.student_id = s.user_id
		JOIN users u ON s.user_id = u.id
		WHERE e.student_id = $1
		ORDER BY e.enrolled_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanEnrollments(rows, total)
}

func (r *repository) IsEnrolled(ctx context.Context, courseID, studentID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM enrollments WHERE course_id = $1 AND student_id = $2 AND status = 'active')`
	var exists bool
	err := r.db.QueryRow(ctx, query, courseID, studentID).Scan(&exists)
	return exists, err
}

func (r *repository) scanEnrollments(rows pgx.Rows, total int64) ([]EnrollmentWithDetails, int64, error) {
	var enrollments []EnrollmentWithDetails
	for rows.Next() {
		var e EnrollmentWithDetails
		if err := rows.Scan(
			&e.ID, &e.CourseID, &e.StudentID, &e.Status, &e.EnrolledAt,
			&e.CourseTitle,
			&e.StudentFirstName, &e.StudentLastName, &e.StudentEmail,
		); err != nil {
			return nil, 0, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, total, nil
}
