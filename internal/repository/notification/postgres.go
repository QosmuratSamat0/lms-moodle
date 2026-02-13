package notification

import (
	"context"

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
		`INSERT INTO notifications (id, user_id, type, title, message, is_read, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		n.ID, n.UserID, n.Type, n.Title, n.Message, n.Read, n.CreatedAt)
	return err
}

func (r *PostgresRepository) CreateBulk(notifications []*notification.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for _, n := range notifications {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO notifications (id, user_id, type, title, message, is_read, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			n.ID, n.UserID, n.Type, n.Title, n.Message, n.Read, n.CreatedAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit(context.Background())
}

func (r *PostgresRepository) GetByID(id string) (*notification.Notification, error) {
	n := &notification.Notification{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, COALESCE(type, 'info'), title, COALESCE(message, ''), is_read, created_at FROM notifications WHERE id = $1`, id).
		Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Read, &n.CreatedAt)
	return n, err
}

func (r *PostgresRepository) ListByUser(userID string, skip, take int) ([]*notification.Notification, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, COALESCE(type, 'info'), title, COALESCE(message, ''), is_read, created_at FROM notifications
		 WHERE user_id=$1 ORDER BY created_at DESC OFFSET $2 LIMIT $3`, userID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*notification.Notification
	for rows.Next() {
		n := &notification.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Read, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (r *PostgresRepository) MarkAsRead(id string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE notifications SET is_read=true WHERE id=$1`, id)
	return err
}

func (r *PostgresRepository) MarkAllAsRead(userID string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE notifications SET is_read=true WHERE user_id=$1 AND is_read=false`, userID)
	return err
}

func (r *PostgresRepository) UnreadCount(userID string) (int, error) {
	var count int
	err := r.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND is_read=false`, userID).
		Scan(&count)
	return count, err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM notifications WHERE id = $1", id)
	return err
}
