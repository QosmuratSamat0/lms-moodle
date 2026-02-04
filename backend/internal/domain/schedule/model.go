package schedule

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents schedule event type
type EventType string

const (
	EventTypeClass     EventType = "class"
	EventTypeLab       EventType = "lab"
	EventTypeExam      EventType = "exam"
	EventTypeOfficeHrs EventType = "office_hours"
	EventTypeEvent     EventType = "event"
)

// RecurrenceType represents recurrence pattern
type RecurrenceType string

const (
	RecurrenceNone    RecurrenceType = "none"
	RecurrenceDaily   RecurrenceType = "daily"
	RecurrenceWeekly  RecurrenceType = "weekly"
	RecurrenceMonthly RecurrenceType = "monthly"
)

// Event represents a schedule event
type Event struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	CourseID    *uuid.UUID     `json:"course_id" db:"course_id"`
	CreatedBy   uuid.UUID      `json:"created_by" db:"created_by"`
	Title       string         `json:"title" db:"title"`
	Description *string        `json:"description" db:"description"`
	EventType   EventType      `json:"event_type" db:"event_type"`
	StartTime   time.Time      `json:"start_time" db:"start_time"`
	EndTime     time.Time      `json:"end_time" db:"end_time"`
	Location    *string        `json:"location" db:"location"`
	IsOnline    bool           `json:"is_online" db:"is_online"`
	MeetingURL  *string        `json:"meeting_url" db:"meeting_url"`
	Recurrence  RecurrenceType `json:"recurrence" db:"recurrence"`
	RecurUntil  *time.Time     `json:"recur_until" db:"recur_until"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// EventWithDetails includes course info
type EventWithDetails struct {
	Event
	CourseTitle  *string `json:"course_title" db:"course_title"`
	CreatorName  string  `json:"creator_name" db:"creator_name"`
	CreatorEmail string  `json:"creator_email" db:"creator_email"`
}
