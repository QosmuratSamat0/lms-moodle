package schedule

import "github.com/google/uuid"

// CreateEventRequest represents request to create a schedule event
type CreateEventRequest struct {
	CourseID    *uuid.UUID `json:"course_id,omitempty"`
	Title       string     `json:"title" validate:"required,min=2,max=200"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	EventType   string     `json:"event_type" validate:"required,oneof=class lab exam office_hours event"`
	StartTime   string     `json:"start_time" validate:"required"`
	EndTime     string     `json:"end_time" validate:"required"`
	Location    *string    `json:"location,omitempty" validate:"omitempty,max=200"`
	IsOnline    bool       `json:"is_online"`
	MeetingURL  *string    `json:"meeting_url,omitempty" validate:"omitempty,url"`
	Recurrence  string     `json:"recurrence" validate:"omitempty,oneof=none daily weekly monthly"`
	RecurUntil  *string    `json:"recur_until,omitempty"`
}

// UpdateEventRequest represents request to update an event
type UpdateEventRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=2,max=200"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	EventType   *string `json:"event_type,omitempty" validate:"omitempty,oneof=class lab exam office_hours event"`
	StartTime   *string `json:"start_time,omitempty"`
	EndTime     *string `json:"end_time,omitempty"`
	Location    *string `json:"location,omitempty" validate:"omitempty,max=200"`
	IsOnline    *bool   `json:"is_online,omitempty"`
	MeetingURL  *string `json:"meeting_url,omitempty" validate:"omitempty,url"`
	Recurrence  *string `json:"recurrence,omitempty" validate:"omitempty,oneof=none daily weekly monthly"`
	RecurUntil  *string `json:"recur_until,omitempty"`
}

// EventResponse represents an event in API responses
type EventResponse struct {
	ID           string  `json:"id"`
	CourseID     *string `json:"course_id,omitempty"`
	CourseTitle  *string `json:"course_title,omitempty"`
	Title        string  `json:"title"`
	Description  *string `json:"description,omitempty"`
	EventType    string  `json:"event_type"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	Location     *string `json:"location,omitempty"`
	IsOnline     bool    `json:"is_online"`
	MeetingURL   *string `json:"meeting_url,omitempty"`
	Recurrence   string  `json:"recurrence"`
	RecurUntil   *string `json:"recur_until,omitempty"`
	CreatorName  string  `json:"creator_name"`
	CreatorEmail string  `json:"creator_email"`
	CreatedAt    string  `json:"created_at"`
}

// EventListResponse represents paginated list of events
type EventListResponse struct {
	Events     []EventResponse `json:"events"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// EventFilter for filtering events
type EventFilter struct {
	CourseID  *uuid.UUID
	StartDate *string
	EndDate   *string
	EventType *string
}
