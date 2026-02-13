package chat

import "time"

// ── Room ────────────────────────────────────────────────────

type Room struct {
	ID        string    `json:"id"        db:"id"`
	CourseID  *string   `json:"course_id" db:"course_id"`
	RoomType  string    `json:"type"      db:"room_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Aggregated fields (filled by queries, not direct columns)
	Name            *string        `json:"name,omitempty"`
	CourseTitle      *string       `json:"course_title,omitempty"`
	MemberCount     int            `json:"member_count"`
	Participants    []Participant  `json:"participants,omitempty"`
	LastMessage     *LastMessage   `json:"last_message,omitempty"`
	LastMessageAt   *time.Time     `json:"last_message_at,omitempty"`
	UnreadCount     int            `json:"unread_count"`
}

type Participant struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type LastMessage struct {
	ID              string `json:"id"`
	Content         string `json:"content"`
	SenderID        string `json:"sender_id"`
	SenderFirstName string `json:"sender_first_name,omitempty"`
	CreatedAt       string `json:"created_at"`
}

// ── Member ──────────────────────────────────────────────────

type Member struct {
	RoomID   string    `json:"room_id"   db:"room_id"`
	UserID   string    `json:"user_id"   db:"user_id"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`

	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     string  `json:"email"`
}

// ── Message ─────────────────────────────────────────────────

type Message struct {
	ID             string    `json:"id"             db:"id"`
	RoomID         string    `json:"room_id"        db:"room_id"`
	SenderUserID   string    `json:"sender_id"      db:"sender_user_id"`
	Content        string    `json:"content"        db:"content"`
	CreatedAt      time.Time `json:"created_at"     db:"created_at"`

	// Joined fields
	SenderFirstName string `json:"sender_first_name,omitempty"`
	SenderLastName  string `json:"sender_last_name,omitempty"`
	SenderEmail     string `json:"sender_email,omitempty"`
}

// ── Inputs ──────────────────────────────────────────────────

type CreateRoomInput struct {
	RoomType       string   // "direct", "group", "course"
	CourseID       *string
	Name           *string
	ParticipantIDs []string // user IDs to add (including creator for groups)
	CreatorUserID  string
}

type CreateMessageInput struct {
	RoomID       string
	SenderUserID string
	Content      string
}

// ── Repository ──────────────────────────────────────────────

type Repository interface {
	// Rooms
	CreateRoom(room *Room) error
	GetRoomByID(id string, currentUserID string) (*Room, error)
	ListRoomsByUser(userID string) ([]*Room, error)
	DeleteRoom(id string) error
	FindDirectRoom(userA, userB string) (*Room, error)

	// Members
	AddMember(roomID, userID string) error
	RemoveMember(roomID, userID string) error
	ListMembers(roomID string) ([]*Member, error)
	IsMember(roomID, userID string) (bool, error)
	GetMemberCount(roomID string) (int, error)

	// Messages
	CreateMessage(msg *Message) error
	GetMessageByID(id string) (*Message, error)
	ListMessages(roomID string, limit, offset int) ([]*Message, int, error)
	DeleteMessage(id string) error

	// Users search
	SearchUsers(query string, limit int) ([]Participant, error)
}
