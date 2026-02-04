package chat

import "time"

type Message struct {
	ID        string    `db:"id"`
	SenderID  string    `db:"sender_id"`
	CourseID  string    `db:"course_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

type CreateMessageInput struct {
	SenderID string
	CourseID string
	Content  string
}

type Repository interface {
	Create(message *Message) error
	GetByID(id string) (*Message, error)
	ListByCourse(courseID string, skip, take int) ([]*Message, error)
	Delete(id string) error
}
