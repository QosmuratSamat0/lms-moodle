package chat

import (
	"fmt"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/chat"
	"github.com/google/uuid"
)

type Service struct {
	repo chat.Repository
}

func NewService(repo chat.Repository) *Service {
	return &Service{repo: repo}
}

// ─── rooms ──────────────────────────────────────────────────

func (s *Service) CreateRoom(input *chat.CreateRoomInput) (*chat.Room, error) {
	room := &chat.Room{
		ID:        uuid.New().String(),
		RoomType:  input.RoomType,
		CourseID:  input.CourseID,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateRoom(room); err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}

	// Add creator as member
	if input.CreatorUserID != "" {
		_ = s.repo.AddMember(room.ID, input.CreatorUserID)
	}

	// Add other participants
	for _, uid := range input.ParticipantIDs {
		if uid != input.CreatorUserID {
			_ = s.repo.AddMember(room.ID, uid)
		}
	}

	return s.repo.GetRoomByID(room.ID, input.CreatorUserID)
}

func (s *Service) GetRoom(id string, currentUserID string) (*chat.Room, error) {
	return s.repo.GetRoomByID(id, currentUserID)
}

func (s *Service) ListRooms(userID string) ([]*chat.Room, error) {
	return s.repo.ListRoomsByUser(userID)
}

func (s *Service) DeleteRoom(id string) error {
	return s.repo.DeleteRoom(id)
}

func (s *Service) GetOrCreateDM(currentUserID, otherUserID string) (*chat.Room, error) {
	// Try to find existing direct room
	room, err := s.repo.FindDirectRoom(currentUserID, otherUserID)
	if err == nil && room != nil {
		return room, nil
	}

	// Create new direct room
	return s.CreateRoom(&chat.CreateRoomInput{
		RoomType:       "direct",
		CreatorUserID:  currentUserID,
		ParticipantIDs: []string{otherUserID},
	})
}

// ─── members ────────────────────────────────────────────────

func (s *Service) AddMember(roomID, userID string) error {
	return s.repo.AddMember(roomID, userID)
}

func (s *Service) RemoveMember(roomID, userID string) error {
	return s.repo.RemoveMember(roomID, userID)
}

func (s *Service) ListMembers(roomID string) ([]*chat.Member, error) {
	return s.repo.ListMembers(roomID)
}

func (s *Service) IsMember(roomID, userID string) (bool, error) {
	return s.repo.IsMember(roomID, userID)
}

// ─── messages ───────────────────────────────────────────────

func (s *Service) SendMessage(input *chat.CreateMessageInput) (*chat.Message, error) {
	msg := &chat.Message{
		ID:           uuid.New().String(),
		RoomID:       input.RoomID,
		SenderUserID: input.SenderUserID,
		Content:      input.Content,
		CreatedAt:    time.Now(),
	}
	if err := s.repo.CreateMessage(msg); err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	// Re-fetch to get joined sender info
	return s.repo.GetMessageByID(msg.ID)
}

func (s *Service) GetMessage(id string) (*chat.Message, error) {
	return s.repo.GetMessageByID(id)
}

func (s *Service) ListMessages(roomID string, limit, offset int) ([]*chat.Message, int, error) {
	return s.repo.ListMessages(roomID, limit, offset)
}

func (s *Service) DeleteMessage(id string) error {
	return s.repo.DeleteMessage(id)
}

// ─── users search ───────────────────────────────────────────

func (s *Service) SearchUsers(query string, limit int) ([]chat.Participant, error) {
	return s.repo.SearchUsers(query, limit)
}
