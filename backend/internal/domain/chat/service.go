package chat

import (
	"context"
	"errors"
	"time"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/ap1-final-mini-moodle/internal/shared/websocket"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the chat service interface
type Service interface {
	// Room operations
	CreateRoom(ctx context.Context, userID uuid.UUID, req *CreateRoomRequest) (*RoomResponse, error)
	GetRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (*RoomResponse, error)
	UpdateRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, req *UpdateRoomRequest) error
	DeleteRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
	ListMyRooms(ctx context.Context, userID uuid.UUID, page, limit int) (*RoomListResponse, error)
	GetOrCreateDirectRoom(ctx context.Context, userID, targetUserID uuid.UUID) (*RoomResponse, error)

	// Member operations
	AddMember(ctx context.Context, roomID uuid.UUID, adderID uuid.UUID, req *AddMemberRequest) error
	RemoveMember(ctx context.Context, roomID uuid.UUID, removerID uuid.UUID, targetUserID uuid.UUID) error
	LeaveRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
	ListMembers(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) ([]MemberResponse, error)

	// Message operations
	SendMessage(ctx context.Context, roomID uuid.UUID, senderID uuid.UUID, req *SendMessageRequest) (*MessageResponse, error)
	UpdateMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, req *UpdateMessageRequest) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error
	ListMessages(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, page, limit int) (*MessageListResponse, error)

	// WebSocket operations
	AuthorizeRoomAccess(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, role string) error
	SaveAndBroadcastMessage(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, content string, hub *websocket.Hub) (*websocket.OutboundMessage, error)
}

type service struct {
	repo Repository
}

// NewService creates a new chat service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Room operations

func (s *service) CreateRoom(ctx context.Context, userID uuid.UUID, req *CreateRoomRequest) (*RoomResponse, error) {
	room := &Room{
		Name:      req.Name,
		Type:      RoomType(req.Type),
		CourseID:  req.CourseID,
		CreatedBy: userID,
	}

	if err := s.repo.CreateRoom(ctx, room); err != nil {
		return nil, errorx.Wrap(err, "create room")
	}

	// Add creator as admin
	member := &Member{
		RoomID: room.ID,
		UserID: userID,
		Role:   "admin",
	}
	if err := s.repo.AddMember(ctx, member); err != nil {
		return nil, errorx.Wrap(err, "add creator as member")
	}

	// Add other members if specified
	for _, memberID := range req.Members {
		if memberID == userID {
			continue
		}
		m := &Member{
			RoomID: room.ID,
			UserID: memberID,
			Role:   "member",
		}
		_ = s.repo.AddMember(ctx, m) // Ignore errors for member addition
	}

	detailed, err := s.repo.GetRoomByID(ctx, room.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created room")
	}

	return s.toRoomResponse(detailed), nil
}

func (s *service) GetRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (*RoomResponse, error) {
	// Check membership
	isMember, err := s.repo.IsMember(ctx, roomID, userID)
	if err != nil {
		return nil, errorx.Wrap(err, "check membership")
	}
	if !isMember {
		return nil, errorx.NewForbiddenError("you are not a member of this room")
	}

	room, err := s.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("chat room")
		}
		return nil, errorx.Wrap(err, "get room")
	}

	return s.toRoomResponse(room), nil
}

func (s *service) UpdateRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, req *UpdateRoomRequest) error {
	room, err := s.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("chat room")
		}
		return errorx.Wrap(err, "get room")
	}

	// Check if user is admin of the room
	member, err := s.repo.GetMember(ctx, roomID, userID)
	if err != nil {
		return errorx.NewForbiddenError("you are not a member of this room")
	}
	if member.Role != "admin" {
		return errorx.NewForbiddenError("only room admins can update the room")
	}

	if req.Name != nil {
		room.Name = req.Name
	}

	if err := s.repo.UpdateRoom(ctx, &room.Room); err != nil {
		return errorx.Wrap(err, "update room")
	}

	return nil
}

func (s *service) DeleteRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	room, err := s.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("chat room")
		}
		return errorx.Wrap(err, "get room")
	}

	// Only creator can delete
	if room.CreatedBy != userID {
		return errorx.NewForbiddenError("only room creator can delete the room")
	}

	if err := s.repo.DeleteRoom(ctx, roomID); err != nil {
		return errorx.Wrap(err, "delete room")
	}

	return nil
}

func (s *service) ListMyRooms(ctx context.Context, userID uuid.UUID, page, limit int) (*RoomListResponse, error) {
	offset := (page - 1) * limit
	rooms, total, err := s.repo.ListUserRooms(ctx, userID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list rooms")
	}

	var responses []RoomResponse
	for _, r := range rooms {
		responses = append(responses, *s.toRoomResponse(&r))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &RoomListResponse{
		Rooms:      responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *service) GetOrCreateDirectRoom(ctx context.Context, userID, targetUserID uuid.UUID) (*RoomResponse, error) {
	// Check if direct room already exists
	existing, err := s.repo.GetDirectRoom(ctx, userID, targetUserID)
	if err == nil {
		detailed, err := s.repo.GetRoomByID(ctx, existing.ID)
		if err != nil {
			return nil, errorx.Wrap(err, "get existing room")
		}
		return s.toRoomResponse(detailed), nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.Wrap(err, "check existing room")
	}

	// Create new direct room
	req := &CreateRoomRequest{
		Type:    "direct",
		Members: []uuid.UUID{targetUserID},
	}

	return s.CreateRoom(ctx, userID, req)
}

// Member operations

func (s *service) AddMember(ctx context.Context, roomID uuid.UUID, adderID uuid.UUID, req *AddMemberRequest) error {
	// Check if adder is admin
	adder, err := s.repo.GetMember(ctx, roomID, adderID)
	if err != nil {
		return errorx.NewForbiddenError("you are not a member of this room")
	}
	if adder.Role != "admin" {
		return errorx.NewForbiddenError("only admins can add members")
	}

	role := "member"
	if req.Role != "" {
		role = req.Role
	}

	member := &Member{
		RoomID: roomID,
		UserID: req.UserID,
		Role:   role,
	}

	if err := s.repo.AddMember(ctx, member); err != nil {
		return errorx.Wrap(err, "add member")
	}

	return nil
}

func (s *service) RemoveMember(ctx context.Context, roomID uuid.UUID, removerID uuid.UUID, targetUserID uuid.UUID) error {
	// Check if remover is admin
	remover, err := s.repo.GetMember(ctx, roomID, removerID)
	if err != nil {
		return errorx.NewForbiddenError("you are not a member of this room")
	}
	if remover.Role != "admin" {
		return errorx.NewForbiddenError("only admins can remove members")
	}

	// Can't remove yourself through this method
	if targetUserID == removerID {
		return errorx.NewBadRequestError("use leave room to remove yourself")
	}

	if err := s.repo.RemoveMember(ctx, roomID, targetUserID); err != nil {
		return errorx.Wrap(err, "remove member")
	}

	return nil
}

func (s *service) LeaveRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	if err := s.repo.RemoveMember(ctx, roomID, userID); err != nil {
		return errorx.Wrap(err, "leave room")
	}
	return nil
}

func (s *service) ListMembers(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) ([]MemberResponse, error) {
	// Check membership
	isMember, err := s.repo.IsMember(ctx, roomID, userID)
	if err != nil {
		return nil, errorx.Wrap(err, "check membership")
	}
	if !isMember {
		return nil, errorx.NewForbiddenError("you are not a member of this room")
	}

	members, err := s.repo.ListRoomMembers(ctx, roomID)
	if err != nil {
		return nil, errorx.Wrap(err, "list members")
	}

	return s.toMemberResponses(members), nil
}

// Message operations

func (s *service) SendMessage(ctx context.Context, roomID uuid.UUID, senderID uuid.UUID, req *SendMessageRequest) (*MessageResponse, error) {
	// Check membership
	isMember, err := s.repo.IsMember(ctx, roomID, senderID)
	if err != nil {
		return nil, errorx.Wrap(err, "check membership")
	}
	if !isMember {
		return nil, errorx.NewForbiddenError("you are not a member of this room")
	}

	message := &Message{
		RoomID:    roomID,
		SenderID:  senderID,
		Content:   req.Content,
		ReplyToID: req.ReplyToID,
	}

	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, errorx.Wrap(err, "create message")
	}

	detailed, err := s.repo.GetMessageByID(ctx, message.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created message")
	}

	return s.toMessageResponse(detailed), nil
}

func (s *service) UpdateMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, req *UpdateMessageRequest) error {
	message, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("message")
		}
		return errorx.Wrap(err, "get message")
	}

	if message.SenderID != userID {
		return errorx.NewForbiddenError("you can only edit your own messages")
	}

	if message.DeletedAt != nil {
		return errorx.NewBadRequestError("cannot edit deleted message")
	}

	message.Content = req.Content
	if err := s.repo.UpdateMessage(ctx, &message.Message); err != nil {
		return errorx.Wrap(err, "update message")
	}

	return nil
}

func (s *service) DeleteMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error {
	message, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("message")
		}
		return errorx.Wrap(err, "get message")
	}

	if message.SenderID != userID {
		return errorx.NewForbiddenError("you can only delete your own messages")
	}

	if err := s.repo.SoftDeleteMessage(ctx, messageID); err != nil {
		return errorx.Wrap(err, "delete message")
	}

	return nil
}

func (s *service) ListMessages(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, page, limit int) (*MessageListResponse, error) {
	// Check membership
	isMember, err := s.repo.IsMember(ctx, roomID, userID)
	if err != nil {
		return nil, errorx.Wrap(err, "check membership")
	}
	if !isMember {
		return nil, errorx.NewForbiddenError("you are not a member of this room")
	}

	offset := (page - 1) * limit
	messages, total, err := s.repo.ListRoomMessages(ctx, roomID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list messages")
	}

	var responses []MessageResponse
	for _, m := range messages {
		responses = append(responses, *s.toMessageResponse(&m))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &MessageListResponse{
		Messages:   responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Helper methods

func (s *service) toRoomResponse(r *RoomWithDetails) *RoomResponse {
	resp := &RoomResponse{
		ID:          r.ID.String(),
		Name:        r.Name,
		Type:        string(r.Type),
		CourseTitle: r.CourseTitle,
		MemberCount: r.MemberCount,
		LastMessage: r.LastMessage,
		UnreadCount: r.UnreadCount,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
	}
	if r.CourseID != nil {
		id := r.CourseID.String()
		resp.CourseID = &id
	}
	if r.LastMessageAt != nil {
		t := r.LastMessageAt.Format(time.RFC3339)
		resp.LastMessageAt = &t
	}
	return resp
}

func (s *service) toMemberResponses(members []MemberWithDetails) []MemberResponse {
	var responses []MemberResponse
	for _, m := range members {
		responses = append(responses, MemberResponse{
			ID:        m.ID.String(),
			UserID:    m.UserID.String(),
			FirstName: m.FirstName,
			LastName:  m.LastName,
			Email:     m.Email,
			Role:      m.Role,
			UserRole:  m.UserRole,
			JoinedAt:  m.JoinedAt.Format(time.RFC3339),
		})
	}
	return responses
}

func (s *service) toMessageResponse(m *MessageWithDetails) *MessageResponse {
	resp := &MessageResponse{
		ID:              m.ID.String(),
		RoomID:          m.RoomID.String(),
		SenderID:        m.SenderID.String(),
		SenderFirstName: m.SenderFirstName,
		SenderLastName:  m.SenderLastName,
		SenderEmail:     m.SenderEmail,
		Content:         m.Content,
		ReplyToContent:  m.ReplyToContent,
		CreatedAt:       m.CreatedAt.Format(time.RFC3339),
	}
	if m.ReplyToID != nil {
		id := m.ReplyToID.String()
		resp.ReplyToID = &id
	}
	if m.EditedAt != nil {
		t := m.EditedAt.Format(time.RFC3339)
		resp.EditedAt = &t
	}
	return resp
}

// WebSocket operations

// AuthorizeRoomAccess checks if a user can access a chat room.
// Returns nil if authorized, error otherwise.
func (s *service) AuthorizeRoomAccess(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, role string) error {
	// Check if room exists
	room, err := s.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("chat room")
		}
		return errorx.Wrap(err, "get room")
	}

	// Check membership
	isMember, err := s.repo.IsMember(ctx, roomID, userID)
	if err != nil {
		return errorx.Wrap(err, "check membership")
	}

	if isMember {
		return nil
	}

	// For admins, allow access to any room
	if role == "admin" {
		return nil
	}

	// Room exists but user is not a member
	_ = room // room is checked above
	return errorx.NewForbiddenError("you are not a member of this chat room")
}

// SaveAndBroadcastMessage saves a message to the database and broadcasts it to all clients in the room.
func (s *service) SaveAndBroadcastMessage(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, content string, hub *websocket.Hub) (*websocket.OutboundMessage, error) {
	// Create message in database
	message := &Message{
		RoomID:   roomID,
		SenderID: userID,
		Content:  content,
	}

	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, errorx.Wrap(err, "create message")
	}

	// Create outbound message
	outMsg := websocket.NewChatMessage(
		message.ID.String(),
		roomID.String(),
		userID.String(),
		content,
		message.CreatedAt,
	)

	// Broadcast to all clients in the room
	hub.Broadcast(roomID, outMsg)

	return outMsg, nil
}

var _ Service = (*service)(nil)
