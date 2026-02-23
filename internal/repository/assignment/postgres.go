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
	if a.CreatedByTeacherID != nil && *a.CreatedByTeacherID != "" {
		_, _ = r.db.Exec(context.Background(),
			`INSERT INTO teachers (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`,
			*a.CreatedByTeacherID)
	}
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO assignments (id, course_id, title, description, max_points, due_at, allow_late, created_by_teacher_id, grading_category, weight_percentage, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, (SELECT id FROM teachers WHERE user_id = $8), $9, $10, $11)`,
		a.ID, a.CourseID, a.Title, a.Description, a.MaxPoints, a.DueAt, a.AllowLate, a.CreatedByTeacherID, a.GradingCategory, a.WeightPercentage, a.CreatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*assignment.Assignment, error) {
	a := &assignment.Assignment{}
	err := r.db.QueryRow(context.Background(),
		`SELECT a.id, a.course_id, a.title, a.description, a.max_points, a.due_at, a.allow_late,
		        COALESCE(a.created_by_teacher_id::text, ''), COALESCE(a.grading_category, 'register_midterm'),
		        COALESCE(a.weight_percentage, 0), a.created_at,
		        COALESCE(t.first_name, ''), COALESCE(t.last_name, '')
		 FROM assignments a
		 LEFT JOIN teachers t ON a.created_by_teacher_id = t.id
		 WHERE a.id = $1`, id).
		Scan(&a.ID, &a.CourseID, &a.Title, &a.Description, &a.MaxPoints, &a.DueAt, &a.AllowLate,
			&a.CreatedByTeacherID, &a.GradingCategory, &a.WeightPercentage, &a.CreatedAt,
			&a.TeacherFirstName, &a.TeacherLastName)
	return a, err
}

func (r *PostgresRepository) ListByCourse(courseID string, skip, take int) ([]*assignment.Assignment, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT a.id, a.course_id, a.title, COALESCE(a.description, ''), a.max_points, a.due_at, a.allow_late,
		        COALESCE(a.grading_category, 'register_midterm'), COALESCE(a.weight_percentage, 0), a.created_at,
		        COALESCE(c.title, '') as course_title,
		        COALESCE(t.first_name, ''), COALESCE(t.last_name, '')
		 FROM assignments a
		 LEFT JOIN courses c ON a.course_id = c.id
		 LEFT JOIN teachers t ON a.created_by_teacher_id = t.id
		 WHERE a.course_id=$1 ORDER BY a.created_at DESC OFFSET $2 LIMIT $3`, courseID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*assignment.Assignment
	for rows.Next() {
		a := &assignment.Assignment{}
		if err := rows.Scan(&a.ID, &a.CourseID, &a.Title, &a.Description, &a.MaxPoints, &a.DueAt, &a.AllowLate,
			&a.GradingCategory, &a.WeightPercentage, &a.CreatedAt, &a.CourseTitle,
			&a.TeacherFirstName, &a.TeacherLastName); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}

func (r *PostgresRepository) Update(a *assignment.Assignment) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE assignments SET title=$1, description=$2, max_points=$3, due_at=$4, allow_late=$5,
		        grading_category=$6, weight_percentage=$7 WHERE id=$8`,
		a.Title, a.Description, a.MaxPoints, a.DueAt, a.AllowLate,
		a.GradingCategory, a.WeightPercentage, a.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM assignments WHERE id = $1", id)
	return err
}
