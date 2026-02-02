//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/chat"
	"github.com/MaqsattoTeam/aLMS/golang-service/tests/testutil"
	"github.com/google/uuid"
)

func TestChatRepository_CreateRoom(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup - clean and seed basic data
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestUserID1, "user1@test.com", "$2a$10$hash", "student")

	repo := chat.NewRepository(pool)

	tests := []struct {
		name    string
		room    *chat.Room
		wantErr bool
	}{
		{
			name: "create group room",
			room: &chat.Room{
				Name:      ptrString("Test Group"),
				Type:      chat.RoomTypeGroup,
				CreatedBy: testutil.TestUserID1,
			},
			wantErr: false,
		},
		{
			name: "create direct room",
			room: &chat.Room{
				Type:      chat.RoomTypeDirect,
				CreatedBy: testutil.TestUserID1,
			},
			wantErr: false,
		},
		{
			name: "create course room",
			room: &chat.Room{
				Name:      ptrString("Course Chat"),
				Type:      chat.RoomTypeCourse,
				CreatedBy: testutil.TestUserID1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateRoom(ctx, tt.room)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRoom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.room.ID == uuid.Nil {
					t.Error("CreateRoom() did not set room ID")
				}
				if tt.room.CreatedAt.IsZero() {
					t.Error("CreateRoom() did not set CreatedAt")
				}
			}
		})
	}
}

func TestChatRepository_IsMember(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestUserID1, "user1@test.com", "$2a$10$hash", "student")
	fixtures.CreateUser(ctx, t, testutil.TestUserID2, "user2@test.com", "$2a$10$hash", "student")

	repo := chat.NewRepository(pool)

	// Create room
	room := &chat.Room{
		Name:      ptrString("Test Room"),
		Type:      chat.RoomTypeGroup,
		CreatedBy: testutil.TestUserID1,
	}
	if err := repo.CreateRoom(ctx, room); err != nil {
		t.Fatalf("failed to create test room: %v", err)
	}

	// Add user1 as member
	member := &chat.Member{
		RoomID: room.ID,
		UserID: testutil.TestUserID1,
		Role:   "admin",
	}
	if err := repo.AddMember(ctx, member); err != nil {
		t.Fatalf("failed to add member: %v", err)
	}

	tests := []struct {
		name       string
		roomID     uuid.UUID
		userID     uuid.UUID
		wantMember bool
		wantErr    bool
	}{
		{
			name:       "user is member",
			roomID:     room.ID,
			userID:     testutil.TestUserID1,
			wantMember: true,
			wantErr:    false,
		},
		{
			name:       "user is not member",
			roomID:     room.ID,
			userID:     testutil.TestUserID2,
			wantMember: false,
			wantErr:    false,
		},
		{
			name:       "non-existent room",
			roomID:     uuid.New(),
			userID:     testutil.TestUserID1,
			wantMember: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isMember, err := repo.IsMember(ctx, tt.roomID, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("IsMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if isMember != tt.wantMember {
				t.Errorf("IsMember() = %v, want %v", isMember, tt.wantMember)
			}
		})
	}
}

func TestChatRepository_CreateMessage(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestUserID1, "user1@test.com", "$2a$10$hash", "student")

	repo := chat.NewRepository(pool)

	// Create room
	room := &chat.Room{
		Name:      ptrString("Test Room"),
		Type:      chat.RoomTypeGroup,
		CreatedBy: testutil.TestUserID1,
	}
	if err := repo.CreateRoom(ctx, room); err != nil {
		t.Fatalf("failed to create test room: %v", err)
	}

	// Add member
	member := &chat.Member{
		RoomID: room.ID,
		UserID: testutil.TestUserID1,
		Role:   "admin",
	}
	if err := repo.AddMember(ctx, member); err != nil {
		t.Fatalf("failed to add member: %v", err)
	}

	tests := []struct {
		name    string
		message *chat.Message
		wantErr bool
	}{
		{
			name: "create simple message",
			message: &chat.Message{
				RoomID:   room.ID,
				SenderID: testutil.TestUserID1,
				Content:  "Hello, World!",
			},
			wantErr: false,
		},
		{
			name: "create message with long content",
			message: &chat.Message{
				RoomID:   room.ID,
				SenderID: testutil.TestUserID1,
				Content:  string(make([]byte, 5000)),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateMessage(ctx, tt.message)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.message.ID == uuid.Nil {
					t.Error("CreateMessage() did not set message ID")
				}
				if tt.message.CreatedAt.IsZero() {
					t.Error("CreateMessage() did not set CreatedAt")
				}
			}
		})
	}
}

func TestChatRepository_ListRoomMessages(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestUserID1, "user1@test.com", "$2a$10$hash", "student")
	fixtures.CreateStudent(ctx, t, testutil.TestUserID1, "Test", "User", "CS-101")

	repo := chat.NewRepository(pool)

	// Create room
	room := &chat.Room{
		Name:      ptrString("Test Room"),
		Type:      chat.RoomTypeGroup,
		CreatedBy: testutil.TestUserID1,
	}
	if err := repo.CreateRoom(ctx, room); err != nil {
		t.Fatalf("failed to create test room: %v", err)
	}

	// Add member
	member := &chat.Member{
		RoomID: room.ID,
		UserID: testutil.TestUserID1,
		Role:   "admin",
	}
	if err := repo.AddMember(ctx, member); err != nil {
		t.Fatalf("failed to add member: %v", err)
	}

	// Create messages
	for i := 0; i < 5; i++ {
		msg := &chat.Message{
			RoomID:   room.ID,
			SenderID: testutil.TestUserID1,
			Content:  "Message " + string(rune('A'+i)),
		}
		if err := repo.CreateMessage(ctx, msg); err != nil {
			t.Fatalf("failed to create message: %v", err)
		}
	}

	tests := []struct {
		name      string
		roomID    uuid.UUID
		limit     int
		offset    int
		wantCount int
		wantTotal int64
	}{
		{
			name:      "get all messages",
			roomID:    room.ID,
			limit:     10,
			offset:    0,
			wantCount: 5,
			wantTotal: 5,
		},
		{
			name:      "get messages with pagination",
			roomID:    room.ID,
			limit:     2,
			offset:    0,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "get messages with offset",
			roomID:    room.ID,
			limit:     10,
			offset:    3,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "non-existent room",
			roomID:    uuid.New(),
			limit:     10,
			offset:    0,
			wantCount: 0,
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages, total, err := repo.ListRoomMessages(ctx, tt.roomID, tt.limit, tt.offset)
			if err != nil {
				t.Errorf("ListRoomMessages() error = %v", err)
				return
			}

			if len(messages) != tt.wantCount {
				t.Errorf("ListRoomMessages() returned %d messages, want %d", len(messages), tt.wantCount)
			}

			if total != tt.wantTotal {
				t.Errorf("ListRoomMessages() total = %d, want %d", total, tt.wantTotal)
			}
		})
	}
}

func TestChatRepository_AddAndRemoveMember(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestUserID1, "user1@test.com", "$2a$10$hash", "student")
	fixtures.CreateUser(ctx, t, testutil.TestUserID2, "user2@test.com", "$2a$10$hash", "student")

	repo := chat.NewRepository(pool)

	// Create room
	room := &chat.Room{
		Name:      ptrString("Test Room"),
		Type:      chat.RoomTypeGroup,
		CreatedBy: testutil.TestUserID1,
	}
	if err := repo.CreateRoom(ctx, room); err != nil {
		t.Fatalf("failed to create test room: %v", err)
	}

	// Test add member
	member := &chat.Member{
		RoomID: room.ID,
		UserID: testutil.TestUserID2,
		Role:   "member",
	}
	if err := repo.AddMember(ctx, member); err != nil {
		t.Fatalf("AddMember() error = %v", err)
	}

	// Verify member was added
	isMember, err := repo.IsMember(ctx, room.ID, testutil.TestUserID2)
	if err != nil {
		t.Fatalf("IsMember() error = %v", err)
	}
	if !isMember {
		t.Error("expected user to be member after AddMember")
	}

	// Verify member ID was set
	if member.ID == uuid.Nil {
		t.Error("AddMember() did not set member ID")
	}

	// Test remove member
	if err := repo.RemoveMember(ctx, room.ID, testutil.TestUserID2); err != nil {
		t.Fatalf("RemoveMember() error = %v", err)
	}

	// Verify member was removed (soft delete via left_at)
	isMember, err = repo.IsMember(ctx, room.ID, testutil.TestUserID2)
	if err != nil {
		t.Fatalf("IsMember() error = %v", err)
	}
	if isMember {
		t.Error("expected user to not be member after RemoveMember")
	}
}

func TestChatRepository_GetRoomByID(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestUserID1, "user1@test.com", "$2a$10$hash", "student")

	repo := chat.NewRepository(pool)

	// Create room
	room := &chat.Room{
		Name:      ptrString("Test Room"),
		Type:      chat.RoomTypeGroup,
		CreatedBy: testutil.TestUserID1,
	}
	if err := repo.CreateRoom(ctx, room); err != nil {
		t.Fatalf("failed to create test room: %v", err)
	}

	tests := []struct {
		name     string
		roomID   uuid.UUID
		wantName string
		wantErr  bool
	}{
		{
			name:     "get existing room",
			roomID:   room.ID,
			wantName: "Test Room",
			wantErr:  false,
		},
		{
			name:    "get non-existing room",
			roomID:  uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetRoomByID(ctx, tt.roomID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoomByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Name == nil || *got.Name != tt.wantName {
					gotName := "<nil>"
					if got.Name != nil {
						gotName = *got.Name
					}
					t.Errorf("GetRoomByID() name = %v, want %v", gotName, tt.wantName)
				}
			}
		})
	}
}

func ptrString(s string) *string {
	return &s
}
