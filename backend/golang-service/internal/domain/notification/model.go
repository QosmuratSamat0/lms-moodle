package notification

import (
	"time"

	"github.com/google/uuid"
)

// Type represents notification type
type Type string

const (
	TypeInfo    Type = "info"
	TypeWarning Type = "warning"
	TypeSuccess Type = "success"
	TypeError   Type = "error"
)

// Channel represents notification delivery channel
type Channel string

const (
	ChannelInApp Channel = "in_app"
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
)

// Notification represents a notification
type Notification struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Type      Type       `json:"type" db:"type"`
	Title     string     `json:"title" db:"title"`
	Message   string     `json:"message" db:"message"`
	Data      *string    `json:"data" db:"data"` // JSON encoded extra data
	Channel   Channel    `json:"channel" db:"channel"`
	ReadAt    *time.Time `json:"read_at" db:"read_at"`
	SentAt    *time.Time `json:"sent_at" db:"sent_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	EmailEnabled bool      `json:"email_enabled" db:"email_enabled"`
	SMSEnabled   bool      `json:"sms_enabled" db:"sms_enabled"`
	PushEnabled  bool      `json:"push_enabled" db:"push_enabled"`
	InAppEnabled bool      `json:"in_app_enabled" db:"in_app_enabled"`
}

// BulkNotification represents a notification to be sent to multiple users
type BulkNotification struct {
	UserIDs []uuid.UUID
	Type    Type
	Title   string
	Message string
	Data    *string
	Channel Channel
}
