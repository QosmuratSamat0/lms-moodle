package notification

import "time"

type Notification struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Title     string     `db:"title"`
	Message   string     `db:"message"`
	Read      bool       `db:"read"`
	CreatedAt time.Time  `db:"created_at"`
	ReadAt    *time.Time `db:"read_at"`
}

type CreateNotificationInput struct {
	UserID  string
	Title   string
	Message string
}

type Repository interface {
	Create(notification *Notification) error
	GetByID(id string) (*Notification, error)
	ListByUser(userID string, skip, take int) ([]*Notification, error)
	MarkAsRead(id string) error
	Delete(id string) error
}
