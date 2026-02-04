package schedule

import (
	"context"
	"errors"
	"time"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the schedule service interface
type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req *CreateEventRequest) (*EventResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*EventResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string, req *UpdateEventRequest) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string) error
	ListByCourse(ctx context.Context, courseID uuid.UUID, filter *EventFilter, page, limit int) (*EventListResponse, error)
	ListMySchedule(ctx context.Context, userID uuid.UUID, filter *EventFilter, page, limit int) (*EventListResponse, error)
	GetCalendarEvents(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]EventResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new schedule service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, req *CreateEventRequest) (*EventResponse, error) {
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, errorx.NewValidationError("start_time: invalid datetime format, use RFC3339")
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, errorx.NewValidationError("end_time: invalid datetime format, use RFC3339")
	}

	if endTime.Before(startTime) {
		return nil, errorx.NewValidationError("end_time: end time must be after start time")
	}

	var recurUntil *time.Time
	if req.RecurUntil != nil {
		ru, err := time.Parse("2006-01-02", *req.RecurUntil)
		if err != nil {
			return nil, errorx.NewValidationError("recur_until: invalid date format, use YYYY-MM-DD")
		}
		recurUntil = &ru
	}

	recurrence := RecurrenceNone
	if req.Recurrence != "" {
		recurrence = RecurrenceType(req.Recurrence)
	}

	event := &Event{
		CourseID:    req.CourseID,
		CreatedBy:   userID,
		Title:       req.Title,
		Description: req.Description,
		EventType:   EventType(req.EventType),
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    req.Location,
		IsOnline:    req.IsOnline,
		MeetingURL:  req.MeetingURL,
		Recurrence:  recurrence,
		RecurUntil:  recurUntil,
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return nil, errorx.Wrap(err, "create event")
	}

	detailed, err := s.repo.GetByID(ctx, event.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created event")
	}

	return s.toResponse(detailed), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*EventResponse, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("schedule event")
		}
		return nil, errorx.Wrap(err, "get event")
	}
	return s.toResponse(event), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string, req *UpdateEventRequest) error {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("schedule event")
		}
		return errorx.Wrap(err, "get event")
	}

	// Check permission: admin or creator
	if role != "admin" && event.CreatedBy != userID {
		return errorx.NewForbiddenError("you can only update events you created")
	}

	// Update fields
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = req.Description
	}
	if req.EventType != nil {
		event.EventType = EventType(*req.EventType)
	}
	if req.StartTime != nil {
		st, err := time.Parse(time.RFC3339, *req.StartTime)
		if err != nil {
			return errorx.NewValidationError("start_time: invalid datetime format")
		}
		event.StartTime = st
	}
	if req.EndTime != nil {
		et, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			return errorx.NewValidationError("end_time: invalid datetime format")
		}
		event.EndTime = et
	}
	if req.Location != nil {
		event.Location = req.Location
	}
	if req.IsOnline != nil {
		event.IsOnline = *req.IsOnline
	}
	if req.MeetingURL != nil {
		event.MeetingURL = req.MeetingURL
	}
	if req.Recurrence != nil {
		event.Recurrence = RecurrenceType(*req.Recurrence)
	}
	if req.RecurUntil != nil {
		ru, err := time.Parse("2006-01-02", *req.RecurUntil)
		if err != nil {
			return errorx.NewValidationError("recur_until: invalid date format")
		}
		event.RecurUntil = &ru
	}

	if err := s.repo.Update(ctx, &event.Event); err != nil {
		return errorx.Wrap(err, "update event")
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string) error {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("schedule event")
		}
		return errorx.Wrap(err, "get event")
	}

	if role != "admin" && event.CreatedBy != userID {
		return errorx.NewForbiddenError("you can only delete events you created")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errorx.Wrap(err, "delete event")
	}

	return nil
}

func (s *service) ListByCourse(ctx context.Context, courseID uuid.UUID, filter *EventFilter, page, limit int) (*EventListResponse, error) {
	offset := (page - 1) * limit
	events, total, err := s.repo.ListByCourse(ctx, courseID, filter, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list course events")
	}

	return s.toListResponse(events, total, page, limit), nil
}

func (s *service) ListMySchedule(ctx context.Context, userID uuid.UUID, filter *EventFilter, page, limit int) (*EventListResponse, error) {
	offset := (page - 1) * limit
	events, total, err := s.repo.ListByUser(ctx, userID, filter, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list user events")
	}

	return s.toListResponse(events, total, page, limit), nil
}

func (s *service) GetCalendarEvents(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]EventResponse, error) {
	events, err := s.repo.ListByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errorx.Wrap(err, "get calendar events")
	}

	var responses []EventResponse
	for _, e := range events {
		responses = append(responses, *s.toResponse(&e))
	}
	return responses, nil
}

// Helper methods

func (s *service) toResponse(e *EventWithDetails) *EventResponse {
	resp := &EventResponse{
		ID:           e.ID.String(),
		CourseTitle:  e.CourseTitle,
		Title:        e.Title,
		Description:  e.Description,
		EventType:    string(e.EventType),
		StartTime:    e.StartTime.Format(time.RFC3339),
		EndTime:      e.EndTime.Format(time.RFC3339),
		Location:     e.Location,
		IsOnline:     e.IsOnline,
		MeetingURL:   e.MeetingURL,
		Recurrence:   string(e.Recurrence),
		CreatorName:  e.CreatorName,
		CreatorEmail: e.CreatorEmail,
		CreatedAt:    e.CreatedAt.Format(time.RFC3339),
	}
	if e.CourseID != nil {
		id := e.CourseID.String()
		resp.CourseID = &id
	}
	if e.RecurUntil != nil {
		ru := e.RecurUntil.Format("2006-01-02")
		resp.RecurUntil = &ru
	}
	return resp
}

func (s *service) toListResponse(events []EventWithDetails, total int64, page, limit int) *EventListResponse {
	var responses []EventResponse
	for _, e := range events {
		responses = append(responses, *s.toResponse(&e))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &EventListResponse{
		Events:     responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

var _ Service = (*service)(nil)
