package chat

import "github.com/google/uuid"

// CreateRoomRequest represents request to create a chat room
type CreateRoomRequest struct {
	Name     *string     `json:"name,omitempty" validate:"omitempty,max=100"`
	Type     string      `json:"type" validate:"required,oneof=direct group course"`
	CourseID *uuid.UUID  `json:"course_id,omitempty"`
	Members  []uuid.UUID `json:"members,omitempty" validate:"dive"`
}

// UpdateRoomRequest represents request to update a room
type UpdateRoomRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,max=100"`
}

// AddMemberRequest represents request to add a member
type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Role   string    `json:"role" validate:"omitempty,oneof=admin member"`
}

// SendMessageRequest represents request to send a message
type SendMessageRequest struct {
	Content   string     `json:"content" validate:"required,min=1,max=5000"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`
}

// UpdateMessageRequest represents request to update a message
type UpdateMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

// RoomResponse represents a room in API responses
type RoomResponse struct {
	ID            string  `json:"id"`
	Name          *string `json:"name,omitempty"`
	Type          string  `json:"type"`
	CourseID      *string `json:"course_id,omitempty"`
	CourseTitle   *string `json:"course_title,omitempty"`
	MemberCount   int     `json:"member_count"`
	LastMessage   *string `json:"last_message,omitempty"`
	LastMessageAt *string `json:"last_message_at,omitempty"`
	UnreadCount   int     `json:"unread_count"`
	CreatedAt     string  `json:"created_at"`
}

// RoomListResponse represents paginated list of rooms
type RoomListResponse struct {
	Rooms      []RoomResponse `json:"rooms"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// MemberResponse represents a member in API responses
type MemberResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	UserRole  string  `json:"user_role"`
	JoinedAt  string  `json:"joined_at"`
}

// MessageResponse represents a message in API responses
type MessageResponse struct {
	ID              string  `json:"id"`
	RoomID          string  `json:"room_id"`
	SenderID        string  `json:"sender_id"`
	SenderFirstName *string `json:"sender_first_name,omitempty"`
	SenderLastName  *string `json:"sender_last_name,omitempty"`
	SenderEmail     string  `json:"sender_email"`
	Content         string  `json:"content"`
	ReplyToID       *string `json:"reply_to_id,omitempty"`
	ReplyToContent  *string `json:"reply_to_content,omitempty"`
	EditedAt        *string `json:"edited_at,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

// MessageListResponse represents paginated list of messages
type MessageListResponse struct {
	Messages   []MessageResponse `json:"messages"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// WebSocket message types
type WSMessage struct {
	Type    string      `json:"type"` // message, typing, read, join, leave
	Payload interface{} `json:"payload"`
}

type WSNewMessage struct {
	RoomID  string          `json:"room_id"`
	Message MessageResponse `json:"message"`
}

type WSTyping struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	IsTyping bool   `json:"is_typing"`
}
