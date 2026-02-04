package upload

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/upload"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) upload.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(u *upload.Upload) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO uploads (id, user_id, file_name, file_url, file_size, mime_type, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		u.ID, u.UserID, u.FileName, u.FileURL, u.FileSize, u.MimeType, u.CreatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*upload.Upload, error) {
	u := &upload.Upload{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, file_name, file_url, file_size, mime_type, created_at FROM uploads WHERE id = $1`, id).
		Scan(&u.ID, &u.UserID, &u.FileName, &u.FileURL, &u.FileSize, &u.MimeType, &u.CreatedAt)
	return u, err
}

func (r *PostgresRepository) ListByUser(userID string, skip, take int) ([]*upload.Upload, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, file_name, file_url, file_size, mime_type, created_at FROM uploads
		 WHERE user_id=$1 ORDER BY created_at DESC OFFSET $2 LIMIT $3`, userID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploads []*upload.Upload
	for rows.Next() {
		u := &upload.Upload{}
		if err := rows.Scan(&u.ID, &u.UserID, &u.FileName, &u.FileURL, &u.FileSize, &u.MimeType, &u.CreatedAt); err != nil {
			return nil, err
		}
		uploads = append(uploads, u)
	}
	return uploads, rows.Err()
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM uploads WHERE id = $1", id)
	return err
}
