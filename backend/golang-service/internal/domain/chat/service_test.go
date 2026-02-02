package chat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/websocket"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// mockRepository is a mock implementation of the Repository interface for testing.
type mockRepository struct {
	getRoomByIDFn     func(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error)
	isMemberFn        func(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
	createMessageFn   func(ctx context.Context, message *Message) error
	getMemberFn       func(ctx context.Context, roomID, userID uuid.UUID) (*Member, error)
	addMemberFn       func(ctx context.Context, member *Member) error
	removeMemberFn    func(ctx context.Context, roomID, userID uuid.UUID) error
	listRoomMembersFn func(ctx context.Context, roomID uuid.UUID) ([]MemberWithDetails, error)
}

func (m *mockRepository) CreateRoom(ctx context.Context, room *Room) error {
	return nil
}

func (m *mockRepository) GetRoomByID(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
	if m.getRoomByIDFn != nil {
		return m.getRoomByIDFn(ctx, id)
	}
	return nil, pgx.ErrNoRows
}

func (m *mockRepository) UpdateRoom(ctx context.Context, room *Room) error {
	return nil
}

func (m *mockRepository) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockRepository) ListUserRooms(ctx context.Context, userID uuid.UUID, limit, offset int) ([]RoomWithDetails, int64, error) {
	return nil, 0, nil
}

func (m *mockRepository) GetDirectRoom(ctx context.Context, user1ID, user2ID uuid.UUID) (*Room, error) {
	return nil, pgx.ErrNoRows
}

func (m *mockRepository) AddMember(ctx context.Context, member *Member) error {
	if m.addMemberFn != nil {
		return m.addMemberFn(ctx, member)
	}
	return nil
}

func (m *mockRepository) RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error {
	if m.removeMemberFn != nil {
		return m.removeMemberFn(ctx, roomID, userID)
	}
	return nil
}

func (m *mockRepository) GetMember(ctx context.Context, roomID, userID uuid.UUID) (*Member, error) {
	if m.getMemberFn != nil {
		return m.getMemberFn(ctx, roomID, userID)
	}
	return nil, pgx.ErrNoRows
}

func (m *mockRepository) ListRoomMembers(ctx context.Context, roomID uuid.UUID) ([]MemberWithDetails, error) {
	if m.listRoomMembersFn != nil {
		return m.listRoomMembersFn(ctx, roomID)
	}
	return nil, nil
}

func (m *mockRepository) IsMember(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	if m.isMemberFn != nil {
		return m.isMemberFn(ctx, roomID, userID)
	}
	return false, nil
}

func (m *mockRepository) CreateMessage(ctx context.Context, message *Message) error {
	if m.createMessageFn != nil {
		return m.createMessageFn(ctx, message)
	}
	message.ID = uuid.New()
	message.CreatedAt = time.Now()
	return nil
}

func (m *mockRepository) GetMessageByID(ctx context.Context, id uuid.UUID) (*MessageWithDetails, error) {
	return nil, pgx.ErrNoRows
}

func (m *mockRepository) UpdateMessage(ctx context.Context, message *Message) error {
	return nil
}

func (m *mockRepository) SoftDeleteMessage(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockRepository) ListRoomMessages(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]MessageWithDetails, int64, error) {
	return nil, 0, nil
}

// Ensure mock implements the interface.
var _ Repository = (*mockRepository)(nil)

func TestService_AuthorizeRoomAccess(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name        string
		roomID      uuid.UUID
		userID      uuid.UUID
		role        string
		setupMock   func(*mockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "member has access",
			roomID: roomID,
			userID: userID,
			role:   "student",
			setupMock: func(m *mockRepository) {
				m.getRoomByIDFn = func(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
					return &RoomWithDetails{
						Room: Room{ID: roomID, Type: RoomTypeGroup},
					}, nil
				}
				m.isMemberFn = func(ctx context.Context, rID, uID uuid.UUID) (bool, error) {
					return true, nil
				}
			},
			wantErr: false,
		},
		{
			name:   "admin has access to any room",
			roomID: roomID,
			userID: userID,
			role:   "admin",
			setupMock: func(m *mockRepository) {
				m.getRoomByIDFn = func(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
					return &RoomWithDetails{
						Room: Room{ID: roomID, Type: RoomTypeGroup},
					}, nil
				}
				m.isMemberFn = func(ctx context.Context, rID, uID uuid.UUID) (bool, error) {
					return false, nil // not a member but admin
				}
			},
			wantErr: false,
		},
		{
			name:   "non-member forbidden",
			roomID: roomID,
			userID: userID,
			role:   "student",
			setupMock: func(m *mockRepository) {
				m.getRoomByIDFn = func(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
					return &RoomWithDetails{
						Room: Room{ID: roomID, Type: RoomTypeGroup},
					}, nil
				}
				m.isMemberFn = func(ctx context.Context, rID, uID uuid.UUID) (bool, error) {
					return false, nil
				}
			},
			wantErr:     true,
			errContains: "not a member",
		},
		{
			name:   "room not found",
			roomID: roomID,
			userID: userID,
			role:   "student",
			setupMock: func(m *mockRepository) {
				m.getRoomByIDFn = func(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
					return nil, pgx.ErrNoRows
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "repository error on get room",
			roomID: roomID,
			userID: userID,
			role:   "student",
			setupMock: func(m *mockRepository) {
				m.getRoomByIDFn = func(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
		{
			name:   "repository error on check membership",
			roomID: roomID,
			userID: userID,
			role:   "student",
			setupMock: func(m *mockRepository) {
				m.getRoomByIDFn = func(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
					return &RoomWithDetails{
						Room: Room{ID: roomID, Type: RoomTypeGroup},
					}, nil
				}
				m.isMemberFn = func(ctx context.Context, rID, uID uuid.UUID) (bool, error) {
					return false, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := NewService(mockRepo)
			err := svc.AuthorizeRoomAccess(context.Background(), tt.roomID, tt.userID, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.errContains != "" && !containsError(err, tt.errContains) {
					t.Errorf("error should contain %q, got %v", tt.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestService_SaveAndBroadcastMessage(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		roomID    uuid.UUID
		userID    uuid.UUID
		content   string
		setupMock func(*mockRepository)
		wantErr   bool
	}{
		{
			name:    "successful save and broadcast",
			roomID:  roomID,
			userID:  userID,
			content: "Hello, world!",
			setupMock: func(m *mockRepository) {
				m.createMessageFn = func(ctx context.Context, message *Message) error {
					message.ID = uuid.New()
					message.CreatedAt = time.Now()
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:    "repository error",
			roomID:  roomID,
			userID:  userID,
			content: "Hello, world!",
			setupMock: func(m *mockRepository) {
				m.createMessageFn = func(ctx context.Context, message *Message) error {
					return errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := NewService(mockRepo)
			hub := websocket.NewHub()
			go hub.Run()

			msg, err := svc.SaveAndBroadcastMessage(context.Background(), tt.roomID, tt.userID, tt.content, hub)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if msg == nil {
					t.Error("expected message, got nil")
				}
				if msg != nil && msg.Content != tt.content {
					t.Errorf("content mismatch: got %q, want %q", msg.Content, tt.content)
				}
			}
		})
	}
}

func TestMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		content string
		valid   bool
	}{
		{
			name:    "valid short message",
			content: "Hello",
			valid:   true,
		},
		{
			name:    "valid max length message",
			content: string(make([]byte, 5000)),
			valid:   true,
		},
		{
			name:    "empty message",
			content: "",
			valid:   false,
		},
		{
			name:    "whitespace only message",
			content: "   ",
			valid:   false,
		},
		{
			name:    "too long message",
			content: string(make([]byte, 5001)),
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateMessageContent(tt.content)
			if valid != tt.valid {
				t.Errorf("validateMessageContent(%q) = %v, want %v", tt.content, valid, tt.valid)
			}
		})
	}
}

// validateMessageContent validates message content according to business rules.
func validateMessageContent(content string) bool {
	trimmed := content
	// Trim whitespace
	for len(trimmed) > 0 && (trimmed[0] == ' ' || trimmed[0] == '\t' || trimmed[0] == '\n' || trimmed[0] == '\r') {
		trimmed = trimmed[1:]
	}
	for len(trimmed) > 0 && (trimmed[len(trimmed)-1] == ' ' || trimmed[len(trimmed)-1] == '\t' || trimmed[len(trimmed)-1] == '\n' || trimmed[len(trimmed)-1] == '\r') {
		trimmed = trimmed[:len(trimmed)-1]
	}
	if trimmed == "" {
		return false
	}
	if len(content) > 5000 {
		return false
	}
	return true
}

// containsError checks if error message contains a substring.
func containsError(err error, substr string) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Check for substring in error message
	for i := 0; i <= len(errStr)-len(substr); i++ {
		if errStr[i:i+len(substr)] == substr {
			return true
		}
	}
	// Also check for AppError types
	var appErr *errorx.AppError
	if errors.As(err, &appErr) {
		msg := appErr.Message
		for i := 0; i <= len(msg)-len(substr); i++ {
			if msg[i:i+len(substr)] == substr {
				return true
			}
		}
	}
	return false
}
