package announcement

import (
	"context"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/announcement"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) announcement.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, a *announcement.Announcement) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO announcements (id, course_id, author_id, title, content, pinned, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		a.ID, a.CourseID, a.AuthorID, a.Title, a.Content, a.Pinned, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*announcement.Announcement, error) {
	a := &announcement.Announcement{}
	err := r.db.QueryRow(ctx,
		`SELECT id, course_id, author_id, title, content, pinned, created_at, updated_at
		 FROM announcements WHERE id = $1`, id).
		Scan(&a.ID, &a.CourseID, &a.AuthorID, &a.Title, &a.Content, &a.Pinned, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *PostgresRepository) GetByCourseID(ctx context.Context, courseID string, limit, offset int) ([]*announcement.Announcement, int64, error) {
	var total int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM announcements WHERE course_id = $1`, courseID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, course_id, author_id, title, content, pinned, created_at, updated_at
		 FROM announcements WHERE course_id = $1
		 ORDER BY pinned DESC, created_at DESC
		 LIMIT $2 OFFSET $3`, courseID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var announcements []*announcement.Announcement
	for rows.Next() {
		a := &announcement.Announcement{}
		if err := rows.Scan(&a.ID, &a.CourseID, &a.AuthorID, &a.Title, &a.Content, &a.Pinned, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, err
		}
		announcements = append(announcements, a)
	}
	return announcements, total, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, a *announcement.Announcement) error {
	a.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		`UPDATE announcements SET title=$1, content=$2, pinned=$3, updated_at=$4 WHERE id=$5`,
		a.Title, a.Content, a.Pinned, a.UpdatedAt, a.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM announcements WHERE id = $1`, id)
	return err
}
