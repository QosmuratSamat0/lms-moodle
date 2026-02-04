package assignment

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/assignment"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) assignment.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(a *assignment.Assignment) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO assignments (id, course_id, title, description, max_points, due_date, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		a.ID, a.CourseID, a.Title, a.Description, a.MaxPoints, a.DueDate, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*assignment.Assignment, error) {
	a := &assignment.Assignment{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, course_id, title, description, max_points, due_date, created_at, updated_at
		 FROM assignments WHERE id = $1`, id).
		Scan(&a.ID, &a.CourseID, &a.Title, &a.Description, &a.MaxPoints, &a.DueDate, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *PostgresRepository) ListByCourse(courseID string, skip, take int) ([]*assignment.Assignment, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, course_id, title, description, max_points, due_date, created_at, updated_at
		 FROM assignments WHERE course_id=$1 OFFSET $2 LIMIT $3`, courseID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*assignment.Assignment
	for rows.Next() {
		a := &assignment.Assignment{}
		if err := rows.Scan(&a.ID, &a.CourseID, &a.Title, &a.Description, &a.MaxPoints, &a.DueDate, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}

func (r *PostgresRepository) Update(a *assignment.Assignment) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE assignments SET title=$1, description=$2, max_points=$3, due_date=$4, updated_at=$5 WHERE id=$6`,
		a.Title, a.Description, a.MaxPoints, a.DueDate, a.UpdatedAt, a.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM assignments WHERE id = $1", id)
	return err
}
