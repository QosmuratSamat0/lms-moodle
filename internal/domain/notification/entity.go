package notification

import "time"

type Notification struct {
	ID        string     `json:"id" db:"id"`
	UserID    string     `json:"user_id" db:"user_id"`
	Type      string     `json:"type" db:"type"`
	Title     string     `json:"title" db:"title"`
	Message   string     `json:"message" db:"message"`
	Read      bool       `json:"is_read" db:"is_read"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type CreateNotificationInput struct {
	UserID  string
	Type    string
	Title   string
	Message string
}

type Repository interface {
	Create(notification *Notification) error
	CreateBulk(notifications []*Notification) error
	GetByID(id string) (*Notification, error)
	ListByUser(userID string, skip, take int) ([]*Notification, error)
	UnreadCount(userID string) (int, error)
	MarkAsRead(id string) error
	MarkAllAsRead(userID string) error
	Delete(id string) error
}
