package chat

import (
	"time"

	"github.com/google/uuid"
)

// RoomType represents the type of chat room
type RoomType string

const (
	RoomTypeDirect RoomType = "direct"
	RoomTypeGroup  RoomType = "group"
	RoomTypeCourse RoomType = "course"
)

// Room represents a chat room
type Room struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      *string    `json:"name" db:"name"`
	Type      RoomType   `json:"type" db:"type"`
	CourseID  *uuid.UUID `json:"course_id" db:"course_id"`
	CreatedBy uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// RoomWithDetails includes additional info
type RoomWithDetails struct {
	Room
	CourseTitle   *string    `json:"course_title" db:"course_title"`
	MemberCount   int        `json:"member_count" db:"member_count"`
	LastMessage   *string    `json:"last_message" db:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at" db:"last_message_at"`
	UnreadCount   int        `json:"unread_count" db:"unread_count"`
}

// Member represents a chat room member
type Member struct {
	ID       uuid.UUID  `json:"id" db:"id"`
	RoomID   uuid.UUID  `json:"room_id" db:"room_id"`
	UserID   uuid.UUID  `json:"user_id" db:"user_id"`
	Role     string     `json:"role" db:"role"` // admin, member
	JoinedAt time.Time  `json:"joined_at" db:"joined_at"`
	LeftAt   *time.Time `json:"left_at" db:"left_at"`
}

// MemberWithDetails includes user info
type MemberWithDetails struct {
	Member
	FirstName *string `json:"first_name" db:"first_name"`
	LastName  *string `json:"last_name" db:"last_name"`
	Email     string  `json:"email" db:"email"`
	UserRole  string  `json:"user_role" db:"user_role"`
}

// Message represents a chat message
type Message struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	RoomID    uuid.UUID  `json:"room_id" db:"room_id"`
	SenderID  uuid.UUID  `json:"sender_id" db:"sender_id"`
	Content   string     `json:"content" db:"content"`
	ReplyToID *uuid.UUID `json:"reply_to_id" db:"reply_to_id"`
	EditedAt  *time.Time `json:"edited_at" db:"edited_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// MessageWithDetails includes sender info
type MessageWithDetails struct {
	Message
	SenderFirstName *string `json:"sender_first_name" db:"sender_first_name"`
	SenderLastName  *string `json:"sender_last_name" db:"sender_last_name"`
	SenderEmail     string  `json:"sender_email" db:"sender_email"`
	ReplyToContent  *string `json:"reply_to_content" db:"reply_to_content"`
}
