package upload

import "time"

type Upload struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	FileName  string    `db:"file_name"`
	FileURL   string    `db:"file_url"`
	FileSize  int64     `db:"file_size"`
	MimeType  string    `db:"mime_type"`
	CreatedAt time.Time `db:"created_at"`
}

type CreateUploadInput struct {
	UserID   string
	FileName string
	FileURL  string
	FileSize int64
	MimeType string
}

type Repository interface {
	Create(upload *Upload) error
	GetByID(id string) (*Upload, error)
	ListByUser(userID string, skip, take int) ([]*Upload, error)
	Delete(id string) error
}
