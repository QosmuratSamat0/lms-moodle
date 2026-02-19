package course

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/course"
	"github.com/ap1-final-mini-moodle/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db    *pgxpool.Pool
	redis *database.RedisClient
}

func NewRepository(db *pgxpool.Pool, redis *database.RedisClient) course.Repository {
	return &PostgresRepository{db: db, redis: redis}
}

func (r *PostgresRepository) Create(c *course.Course) error {
	ctx := context.Background()

	_, err := r.db.Exec(context.Background(),
		`INSERT INTO courses (id, code, title, description, owner_teacher_id, max_points, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		c.ID, c.Code, c.Title, c.Description, c.TeacherID, c.MaxPoints, c.Active, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err == nil {
		_ = r.redis.Set(ctx, "course:"+c.ID, data, 10*time.Minute)
		_ = r.redis.Set(ctx, "course:code:"+c.Code, data, 10*time.Minute)
	}
	return err
}

func (r *PostgresRepository) GetByID(id string) (*course.Course, error) {
	ctx := context.Background()
	cacheKey := "course:" + id

	cached, err := r.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var c course.Course
		if err := json.Unmarshal([]byte(cached), &c); err == nil {
			return &c, nil
		}
	}

	c := &course.Course{}
	err = r.db.QueryRow(context.Background(),
		`SELECT c.id, c.code, c.title, c.description, c.owner_teacher_id, c.max_points, c.is_active, c.created_at, c.updated_at,
		        COALESCE(t.first_name, ''), COALESCE(t.last_name, '')
		 FROM courses c
		 LEFT JOIN teachers t ON t.id = c.owner_teacher_id
		 WHERE c.id = $1`, id).
		Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.TeacherID, &c.MaxPoints, &c.Active, &c.CreatedAt, &c.UpdatedAt,
			&c.TeacherFirstName, &c.TeacherLastName)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(c)
	if err == nil {
		r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return c, err
}

func (r *PostgresRepository) GetByCode(code string) (*course.Course, error) {
	ctx := context.Background()
	cacheKey := "course:code:" + code
	cached, err := r.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var c course.Course
		if err := json.Unmarshal([]byte(cached), &c); err == nil {
			return &c, nil
		}
	}

	c := &course.Course{}
	err = r.db.QueryRow(context.Background(),
		`SELECT c.id, c.code, c.title, c.description, c.owner_teacher_id, c.max_points, c.is_active, c.created_at, c.updated_at,
		        COALESCE(t.first_name, ''), COALESCE(t.last_name, '')
		 FROM courses c
		 LEFT JOIN teachers t ON t.id = c.owner_teacher_id
		 WHERE c.code = $1`, code).
		Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.TeacherID, &c.MaxPoints, &c.Active, &c.CreatedAt, &c.UpdatedAt,
			&c.TeacherFirstName, &c.TeacherLastName)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(c)
	r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	return c, err
}

func (r *PostgresRepository) Update(c *course.Course) error {
	ctx := context.Background()
	_, err := r.db.Exec(context.Background(),
		`UPDATE courses SET title=$1, description=$2, max_points=$3, is_active=$4, updated_at=$5 WHERE id=$6`,
		c.Title, c.Description, c.MaxPoints, c.Active, c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}
	r.redis.Del(ctx, "course:"+c.ID)
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
	if err != nil {
		return err
	}
	r.redis.Del(context.Background(), "course:"+id)
	return err
}
