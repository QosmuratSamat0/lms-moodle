package course

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/course"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) course.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(c *course.Course) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO courses (id, code, title, description, owner_teacher_id, max_points, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		c.ID, c.Code, c.Title, c.Description, c.TeacherID, c.MaxPoints, c.Active, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*course.Course, error) {
	c := &course.Course{}
	err := r.db.QueryRow(context.Background(),
		`SELECT c.id, c.code, c.title, c.description, c.owner_teacher_id, c.max_points, c.is_active, c.created_at, c.updated_at,
		        COALESCE(t.first_name, ''), COALESCE(t.last_name, '')
		 FROM courses c
		 LEFT JOIN teachers t ON t.id = c.owner_teacher_id
		 WHERE c.id = $1`, id).
		Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.TeacherID, &c.MaxPoints, &c.Active, &c.CreatedAt, &c.UpdatedAt,
			&c.TeacherFirstName, &c.TeacherLastName)
	return c, err
}

func (r *PostgresRepository) GetByCode(code string) (*course.Course, error) {
	c := &course.Course{}
	err := r.db.QueryRow(context.Background(),
		`SELECT c.id, c.code, c.title, c.description, c.owner_teacher_id, c.max_points, c.is_active, c.created_at, c.updated_at,
		        COALESCE(t.first_name, ''), COALESCE(t.last_name, '')
		 FROM courses c
		 LEFT JOIN teachers t ON t.id = c.owner_teacher_id
		 WHERE c.code = $1`, code).
		Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.TeacherID, &c.MaxPoints, &c.Active, &c.CreatedAt, &c.UpdatedAt,
			&c.TeacherFirstName, &c.TeacherLastName)
	return c, err
}

func (r *PostgresRepository) Update(c *course.Course) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE courses SET title=$1, description=$2, max_points=$3, is_active=$4, updated_at=$5 WHERE id=$6`,
		c.Title, c.Description, c.MaxPoints, c.Active, c.UpdatedAt, c.ID)
	return err
}

func (r *PostgresRepository) List(skip, take int) ([]*course.Course, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT c.id, c.code, c.title, c.description, c.owner_teacher_id, c.max_points, c.is_active, c.created_at, c.updated_at,
		        COALESCE(t.first_name, ''), COALESCE(t.last_name, '')
		 FROM courses c
		 LEFT JOIN teachers t ON t.id = c.owner_teacher_id
		 WHERE c.is_active=true OFFSET $1 LIMIT $2`, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*course.Course
	for rows.Next() {
		c := &course.Course{}
		if err := rows.Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.TeacherID, &c.MaxPoints, &c.Active, &c.CreatedAt, &c.UpdatedAt,
			&c.TeacherFirstName, &c.TeacherLastName); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

func (r *PostgresRepository) ListByTeacher(teacherID string, skip, take int) ([]*course.Course, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT c.id, c.code, c.title, c.description, c.owner_teacher_id, c.max_points, c.is_active, c.created_at, c.updated_at,
		        COALESCE(t.first_name, ''), COALESCE(t.last_name, '')
		 FROM courses c
		 LEFT JOIN teachers t ON t.id = c.owner_teacher_id
		 WHERE c.owner_teacher_id=$1 OFFSET $2 LIMIT $3`, teacherID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*course.Course
	for rows.Next() {
		c := &course.Course{}
		if err := rows.Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.TeacherID, &c.MaxPoints, &c.Active, &c.CreatedAt, &c.UpdatedAt,
			&c.TeacherFirstName, &c.TeacherLastName); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM courses WHERE id = $1", id)
	return err
}
