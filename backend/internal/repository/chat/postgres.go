package chat

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/chat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) chat.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(m *chat.Message) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO messages (id, sender_id, course_id, content, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		m.ID, m.SenderID, m.CourseID, m.Content, m.CreatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*chat.Message, error) {
	m := &chat.Message{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, sender_id, course_id, content, created_at FROM messages WHERE id = $1`, id).
		Scan(&m.ID, &m.SenderID, &m.CourseID, &m.Content, &m.CreatedAt)
	return m, err
}

func (r *PostgresRepository) ListByCourse(courseID string, skip, take int) ([]*chat.Message, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, sender_id, course_id, content, created_at FROM messages
		 WHERE course_id=$1 ORDER BY created_at DESC OFFSET $2 LIMIT $3`, courseID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*chat.Message
	for rows.Next() {
		m := &chat.Message{}
		if err := rows.Scan(&m.ID, &m.SenderID, &m.CourseID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM messages WHERE id = $1", id)
	return err
}
