package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the notification repository interface
type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	CreateBulk(ctx context.Context, notifications []*Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]Notification, int64, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteAll(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new notification repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, notification *Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, title, message, data, channel)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Data,
		notification.Channel,
	).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *repository) CreateBulk(ctx context.Context, notifications []*Notification) error {
	batch := &pgx.Batch{}

	for _, n := range notifications {
		batch.Queue(`
			INSERT INTO notifications (user_id, type, title, message, data, channel)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			n.UserID, n.Type, n.Title, n.Message, n.Data, n.Channel,
		)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for range notifications {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data, channel, read_at, sent_at, created_at
		FROM notifications
		WHERE id = $1`

	var n Notification
	err := r.db.QueryRow(ctx, query, id).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Data,
		&n.Channel, &n.ReadAt, &n.SentAt, &n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *repository) ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]Notification, int64, error) {
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	listQuery := `
		SELECT id, user_id, type, title, message, data, channel, read_at, sent_at, created_at
		FROM notifications
		WHERE user_id = $1`

	if unreadOnly {
		countQuery += " AND read_at IS NULL"
		listQuery += " AND read_at IS NULL"
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery += " ORDER BY created_at DESC LIMIT $2 OFFSET $3"

	rows, err := r.db.Query(ctx, listQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Data,
			&n.Channel, &n.ReadAt, &n.SentAt, &n.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

func (r *repository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET read_at = now() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET read_at = now() WHERE user_id = $1 AND read_at IS NULL`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) DeleteAll(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM notifications WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *repository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read_at IS NULL`
	var count int64
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}
