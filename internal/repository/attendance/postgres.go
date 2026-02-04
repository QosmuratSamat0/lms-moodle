package attendance

import (
	"context"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/attendance"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) attendance.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(a *attendance.Attendance) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO attendance (id, course_id, student_id, date, present, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		a.ID, a.CourseID, a.StudentID, a.Date, a.Present, a.CreatedAt)
	return err
}

func (r *PostgresRepository) ListByCourse(courseID string, skip, take int) ([]*attendance.Attendance, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, course_id, student_id, date, present, created_at FROM attendance
		 WHERE course_id=$1 OFFSET $2 LIMIT $3`, courseID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*attendance.Attendance
	for rows.Next() {
		a := &attendance.Attendance{}
		if err := rows.Scan(&a.ID, &a.CourseID, &a.StudentID, &a.Date, &a.Present, &a.CreatedAt); err != nil {
			return nil, err
		}
		attendances = append(attendances, a)
	}
	return attendances, rows.Err()
}

func (r *PostgresRepository) ListByStudent(studentID string, skip, take int) ([]*attendance.Attendance, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, course_id, student_id, date, present, created_at FROM attendance
		 WHERE student_id=$1 OFFSET $2 LIMIT $3`, studentID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*attendance.Attendance
	for rows.Next() {
		a := &attendance.Attendance{}
		if err := rows.Scan(&a.ID, &a.CourseID, &a.StudentID, &a.Date, &a.Present, &a.CreatedAt); err != nil {
			return nil, err
		}
		attendances = append(attendances, a)
	}
	return attendances, rows.Err()
}

func (r *PostgresRepository) GetByStudentAndDate(studentID string, date time.Time) (*attendance.Attendance, error) {
	a := &attendance.Attendance{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, course_id, student_id, date, present, created_at FROM attendance
		 WHERE student_id=$1 AND DATE(date)=DATE($2)`, studentID, date).
		Scan(&a.ID, &a.CourseID, &a.StudentID, &a.Date, &a.Present, &a.CreatedAt)
	return a, err
}

func (r *PostgresRepository) Update(a *attendance.Attendance) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE attendance SET present=$1 WHERE id=$2`, a.Present, a.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM attendance WHERE id = $1", id)
	return err
}
