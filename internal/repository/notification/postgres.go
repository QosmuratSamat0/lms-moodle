package notification

import (
	"context"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/notification"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) notification.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(n *notification.Notification) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO notifications (id, user_id, title, message, read, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		n.ID, n.UserID, n.Title, n.Message, n.Read, n.CreatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*notification.Notification, error) {
	n := &notification.Notification{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, title, message, read, created_at, read_at FROM notifications WHERE id = $1`, id).
		Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Read, &n.CreatedAt, &n.ReadAt)
	return n, err
}

func (r *PostgresRepository) ListByUser(userID string, skip, take int) ([]*notification.Notification, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, title, message, read, created_at, read_at FROM notifications
		 WHERE user_id=$1 ORDER BY created_at DESC OFFSET $2 LIMIT $3`, userID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*notification.Notification
	for rows.Next() {
		n := &notification.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Read, &n.CreatedAt, &n.ReadAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (r *PostgresRepository) MarkAsRead(id string) error {
	now := time.Now()
	_, err := r.db.Exec(context.Background(),
		`UPDATE notifications SET read=true, read_at=$1 WHERE id=$2`, now, id)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM notifications WHERE id = $1", id)
	return err
}
