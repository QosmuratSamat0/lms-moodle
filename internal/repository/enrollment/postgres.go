package enrollment

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/enrollment"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) enrollment.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(e *enrollment.Enrollment) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO enrollments (id, course_id, student_id, enrolled_at, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		e.ID, e.CourseID, e.StudentID, e.EnrolledAt, e.Status, e.CreatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*enrollment.Enrollment, error) {
	e := &enrollment.Enrollment{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, course_id, student_id, enrolled_at, status, created_at FROM enrollments WHERE id = $1`, id).
		Scan(&e.ID, &e.CourseID, &e.StudentID, &e.EnrolledAt, &e.Status, &e.CreatedAt)
	return e, err
}

func (r *PostgresRepository) GetByCourseAndStudent(courseID, studentID string) (*enrollment.Enrollment, error) {
	e := &enrollment.Enrollment{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, course_id, student_id, enrolled_at, status, created_at 
		 FROM enrollments WHERE course_id=$1 AND student_id=$2`, courseID, studentID).
		Scan(&e.ID, &e.CourseID, &e.StudentID, &e.EnrolledAt, &e.Status, &e.CreatedAt)
	return e, err
}

func (r *PostgresRepository) ListByCourse(courseID string, skip, take int) ([]*enrollment.Enrollment, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, course_id, student_id, enrolled_at, status, created_at 
		 FROM enrollments WHERE course_id=$1 OFFSET $2 LIMIT $3`, courseID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*enrollment.Enrollment
	for rows.Next() {
		e := &enrollment.Enrollment{}
		if err := rows.Scan(&e.ID, &e.CourseID, &e.StudentID, &e.EnrolledAt, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, rows.Err()
}

func (r *PostgresRepository) ListByStudent(studentID string, skip, take int) ([]*enrollment.Enrollment, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, course_id, student_id, enrolled_at, status, created_at 
		 FROM enrollments WHERE student_id=$1 OFFSET $2 LIMIT $3`, studentID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*enrollment.Enrollment
	for rows.Next() {
		e := &enrollment.Enrollment{}
		if err := rows.Scan(&e.ID, &e.CourseID, &e.StudentID, &e.EnrolledAt, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, rows.Err()
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM enrollments WHERE id = $1", id)
	return err
}
