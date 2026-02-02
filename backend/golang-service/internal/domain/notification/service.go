package notification

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the notification service interface
type Service interface {
	// CRUD operations
	Create(ctx context.Context, userID uuid.UUID, notificationType Type, title, message string, data interface{}) error
	CreateBulk(ctx context.Context, bulk *BulkNotification) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*NotificationResponse, error)
	List(ctx context.Context, userID uuid.UUID, unreadOnly bool, page, limit int) (*NotificationListResponse, error)
	MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	DeleteAll(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)

	// Helper methods for common notifications
	NotifyNewAssignment(ctx context.Context, studentIDs []uuid.UUID, assignmentTitle, courseName string, dueDate time.Time) error
	NotifyGradePosted(ctx context.Context, studentID uuid.UUID, assignmentTitle string, score int, maxPoints int) error
	NotifyEnrollmentApproved(ctx context.Context, studentID uuid.UUID, courseName string) error
	NotifyNewMessage(ctx context.Context, userID uuid.UUID, senderName, roomName string) error
}

type service struct {
	repo Repository
}

// NewService creates a new notification service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, notificationType Type, title, message string, data interface{}) error {
	var dataStr *string
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return errorx.Wrap(err, "marshal notification data")
		}
		str := string(jsonData)
		dataStr = &str
	}

	notification := &Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Message: message,
		Data:    dataStr,
		Channel: ChannelInApp,
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return errorx.Wrap(err, "create notification")
	}

	return nil
}

func (s *service) CreateBulk(ctx context.Context, bulk *BulkNotification) error {
	var notifications []*Notification
	for _, userID := range bulk.UserIDs {
		notifications = append(notifications, &Notification{
			UserID:  userID,
			Type:    bulk.Type,
			Title:   bulk.Title,
			Message: bulk.Message,
			Data:    bulk.Data,
			Channel: bulk.Channel,
		})
	}

	if err := s.repo.CreateBulk(ctx, notifications); err != nil {
		return errorx.Wrap(err, "create bulk notifications")
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*NotificationResponse, error) {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("notification")
		}
		return nil, errorx.Wrap(err, "get notification")
	}

	// Verify ownership
	if notification.UserID != userID {
		return nil, errorx.NewForbiddenError("you can only view your own notifications")
	}

	return s.toResponse(notification), nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID, unreadOnly bool, page, limit int) (*NotificationListResponse, error) {
	offset := (page - 1) * limit
	notifications, total, err := s.repo.ListByUser(ctx, userID, unreadOnly, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list notifications")
	}

	var responses []NotificationResponse
	for _, n := range notifications {
		responses = append(responses, *s.toResponse(&n))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &NotificationListResponse{
		Notifications: responses,
		Total:         total,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
	}, nil
}

func (s *service) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("notification")
		}
		return errorx.Wrap(err, "get notification")
	}

	if notification.UserID != userID {
		return errorx.NewForbiddenError("you can only mark your own notifications as read")
	}

	if err := s.repo.MarkAsRead(ctx, id); err != nil {
		return errorx.Wrap(err, "mark as read")
	}

	return nil
}

func (s *service) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.MarkAllAsRead(ctx, userID); err != nil {
		return errorx.Wrap(err, "mark all as read")
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("notification")
		}
		return errorx.Wrap(err, "get notification")
	}

	if notification.UserID != userID {
		return errorx.NewForbiddenError("you can only delete your own notifications")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errorx.Wrap(err, "delete notification")
	}

	return nil
}

func (s *service) DeleteAll(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.DeleteAll(ctx, userID); err != nil {
		return errorx.Wrap(err, "delete all notifications")
	}
	return nil
}

func (s *service) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return 0, errorx.Wrap(err, "count unread")
	}
	return count, nil
}

// Helper methods for common notifications

func (s *service) NotifyNewAssignment(ctx context.Context, studentIDs []uuid.UUID, assignmentTitle, courseName string, dueDate time.Time) error {
	data := map[string]interface{}{
		"assignment_title": assignmentTitle,
		"course_name":      courseName,
		"due_date":         dueDate.Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(data)
	dataStr := string(jsonData)

	bulk := &BulkNotification{
		UserIDs: studentIDs,
		Type:    TypeInfo,
		Title:   "New Assignment: " + assignmentTitle,
		Message: "A new assignment has been posted in " + courseName + ". Due: " + dueDate.Format("Jan 02, 2006 15:04"),
		Data:    &dataStr,
		Channel: ChannelInApp,
	}

	return s.CreateBulk(ctx, bulk)
}

func (s *service) NotifyGradePosted(ctx context.Context, studentID uuid.UUID, assignmentTitle string, score int, maxPoints int) error {
	data := map[string]interface{}{
		"assignment_title": assignmentTitle,
		"score":            score,
		"max_points":       maxPoints,
	}

	return s.Create(ctx, studentID, TypeSuccess, "Grade Posted",
		"Your submission for \""+assignmentTitle+"\" has been graded: "+string(rune(score))+"/"+string(rune(maxPoints)),
		data,
	)
}

func (s *service) NotifyEnrollmentApproved(ctx context.Context, studentID uuid.UUID, courseName string) error {
	return s.Create(ctx, studentID, TypeSuccess, "Enrollment Approved",
		"You have been enrolled in "+courseName,
		map[string]string{"course_name": courseName},
	)
}

func (s *service) NotifyNewMessage(ctx context.Context, userID uuid.UUID, senderName, roomName string) error {
	return s.Create(ctx, userID, TypeInfo, "New Message",
		senderName+" sent a message in "+roomName,
		map[string]string{"sender_name": senderName, "room_name": roomName},
	)
}

// Helper methods

func (s *service) toResponse(n *Notification) *NotificationResponse {
	resp := &NotificationResponse{
		ID:        n.ID.String(),
		Type:      string(n.Type),
		Title:     n.Title,
		Message:   n.Message,
		Channel:   string(n.Channel),
		IsRead:    n.ReadAt != nil,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
	if n.Data != nil {
		resp.Data = n.Data
	}
	if n.ReadAt != nil {
		t := n.ReadAt.Format(time.RFC3339)
		resp.ReadAt = &t
	}
	return resp
}

var _ Service = (*service)(nil)

// NotificationResponse represents a notification in API responses
type NotificationResponse struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	Message   string  `json:"message"`
	Data      *string `json:"data,omitempty"`
	Channel   string  `json:"channel"`
	IsRead    bool    `json:"is_read"`
	ReadAt    *string `json:"read_at,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// NotificationListResponse represents paginated list of notifications
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}
