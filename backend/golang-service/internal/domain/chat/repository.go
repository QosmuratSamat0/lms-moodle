package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the chat repository interface
type Repository interface {
	// Room methods
	CreateRoom(ctx context.Context, room *Room) error
	GetRoomByID(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error)
	UpdateRoom(ctx context.Context, room *Room) error
	DeleteRoom(ctx context.Context, id uuid.UUID) error
	ListUserRooms(ctx context.Context, userID uuid.UUID, limit, offset int) ([]RoomWithDetails, int64, error)
	GetDirectRoom(ctx context.Context, user1ID, user2ID uuid.UUID) (*Room, error)

	// Member methods
	AddMember(ctx context.Context, member *Member) error
	RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error
	GetMember(ctx context.Context, roomID, userID uuid.UUID) (*Member, error)
	ListRoomMembers(ctx context.Context, roomID uuid.UUID) ([]MemberWithDetails, error)
	IsMember(ctx context.Context, roomID, userID uuid.UUID) (bool, error)

	// Message methods
	CreateMessage(ctx context.Context, message *Message) error
	GetMessageByID(ctx context.Context, id uuid.UUID) (*MessageWithDetails, error)
	UpdateMessage(ctx context.Context, message *Message) error
	SoftDeleteMessage(ctx context.Context, id uuid.UUID) error
	ListRoomMessages(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]MessageWithDetails, int64, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new chat repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

// Room methods

func (r *repository) CreateRoom(ctx context.Context, room *Room) error {
	query := `
		INSERT INTO chat_rooms (name, type, course_id, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		room.Name, room.Type, room.CourseID, room.CreatedBy,
	).Scan(&room.ID, &room.CreatedAt)
}

func (r *repository) GetRoomByID(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
	query := `
		SELECT 
			r.id, r.name, r.type, r.course_id, r.created_by, r.created_at,
			c.title AS course_title,
			(SELECT COUNT(*) FROM chat_room_members WHERE room_id = r.id AND left_at IS NULL) AS member_count,
			(SELECT content FROM chat_messages WHERE room_id = r.id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) AS last_message,
			(SELECT created_at FROM chat_messages WHERE room_id = r.id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) AS last_message_at
		FROM chat_rooms r
		LEFT JOIN courses c ON r.course_id = c.id
		WHERE r.id = $1`

	var room RoomWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&room.ID, &room.Name, &room.Type, &room.CourseID, &room.CreatedBy, &room.CreatedAt,
		&room.CourseTitle, &room.MemberCount, &room.LastMessage, &room.LastMessageAt,
	)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *repository) UpdateRoom(ctx context.Context, room *Room) error {
	query := `UPDATE chat_rooms SET name = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, room.Name, room.ID)
	return err
}

func (r *repository) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	// Delete messages first
	_, _ = r.db.Exec(ctx, "DELETE FROM chat_messages WHERE room_id = $1", id)
	// Delete members
	_, _ = r.db.Exec(ctx, "DELETE FROM chat_room_members WHERE room_id = $1", id)
	// Delete room
	_, err := r.db.Exec(ctx, "DELETE FROM chat_rooms WHERE id = $1", id)
	return err
}

func (r *repository) ListUserRooms(ctx context.Context, userID uuid.UUID, limit, offset int) ([]RoomWithDetails, int64, error) {
	countQuery := `
		SELECT COUNT(*) FROM chat_rooms r
		JOIN chat_room_members m ON r.id = m.room_id
		WHERE m.user_id = $1 AND m.left_at IS NULL`

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			r.id, r.name, r.type, r.course_id, r.created_by, r.created_at,
			c.title AS course_title,
			(SELECT COUNT(*) FROM chat_room_members WHERE room_id = r.id AND left_at IS NULL) AS member_count,
			(SELECT content FROM chat_messages WHERE room_id = r.id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) AS last_message,
			(SELECT created_at FROM chat_messages WHERE room_id = r.id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) AS last_message_at,
			0 AS unread_count
		FROM chat_rooms r
		JOIN chat_room_members m ON r.id = m.room_id
		LEFT JOIN courses c ON r.course_id = c.id
		WHERE m.user_id = $1 AND m.left_at IS NULL
		ORDER BY last_message_at DESC NULLS LAST
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanRooms(rows, total)
}

func (r *repository) GetDirectRoom(ctx context.Context, user1ID, user2ID uuid.UUID) (*Room, error) {
	query := `
		SELECT r.id, r.name, r.type, r.course_id, r.created_by, r.created_at
		FROM chat_rooms r
		WHERE r.type = 'direct'
		AND EXISTS (SELECT 1 FROM chat_room_members WHERE room_id = r.id AND user_id = $1 AND left_at IS NULL)
		AND EXISTS (SELECT 1 FROM chat_room_members WHERE room_id = r.id AND user_id = $2 AND left_at IS NULL)
		AND (SELECT COUNT(*) FROM chat_room_members WHERE room_id = r.id AND left_at IS NULL) = 2`

	var room Room
	err := r.db.QueryRow(ctx, query, user1ID, user2ID).Scan(
		&room.ID, &room.Name, &room.Type, &room.CourseID, &room.CreatedBy, &room.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *repository) scanRooms(rows pgx.Rows, total int64) ([]RoomWithDetails, int64, error) {
	var rooms []RoomWithDetails
	for rows.Next() {
		var room RoomWithDetails
		if err := rows.Scan(
			&room.ID, &room.Name, &room.Type, &room.CourseID, &room.CreatedBy, &room.CreatedAt,
			&room.CourseTitle, &room.MemberCount, &room.LastMessage, &room.LastMessageAt, &room.UnreadCount,
		); err != nil {
			return nil, 0, err
		}
		rooms = append(rooms, room)
	}
	return rooms, total, nil
}

// Member methods

func (r *repository) AddMember(ctx context.Context, member *Member) error {
	query := `
		INSERT INTO chat_room_members (room_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (room_id, user_id) DO UPDATE SET left_at = NULL, role = $3
		RETURNING id, joined_at`

	return r.db.QueryRow(ctx, query, member.RoomID, member.UserID, member.Role).Scan(&member.ID, &member.JoinedAt)
}

func (r *repository) RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error {
	query := `UPDATE chat_room_members SET left_at = now() WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, roomID, userID)
	return err
}

func (r *repository) GetMember(ctx context.Context, roomID, userID uuid.UUID) (*Member, error) {
	query := `
		SELECT id, room_id, user_id, role, joined_at, left_at
		FROM chat_room_members
		WHERE room_id = $1 AND user_id = $2 AND left_at IS NULL`

	var m Member
	err := r.db.QueryRow(ctx, query, roomID, userID).Scan(
		&m.ID, &m.RoomID, &m.UserID, &m.Role, &m.JoinedAt, &m.LeftAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) ListRoomMembers(ctx context.Context, roomID uuid.UUID) ([]MemberWithDetails, error) {
	query := `
		SELECT 
			m.id, m.room_id, m.user_id, m.role, m.joined_at, m.left_at,
			COALESCE(s.first_name, t.first_name, mg.first_name) AS first_name,
			COALESCE(s.last_name, t.last_name, mg.last_name) AS last_name,
			u.email, u.role AS user_role
		FROM chat_room_members m
		JOIN users u ON m.user_id = u.id
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN teachers t ON u.id = t.user_id
		LEFT JOIN managers mg ON u.id = mg.user_id
		WHERE m.room_id = $1 AND m.left_at IS NULL
		ORDER BY m.joined_at`

	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []MemberWithDetails
	for rows.Next() {
		var m MemberWithDetails
		if err := rows.Scan(
			&m.ID, &m.RoomID, &m.UserID, &m.Role, &m.JoinedAt, &m.LeftAt,
			&m.FirstName, &m.LastName, &m.Email, &m.UserRole,
		); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *repository) IsMember(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM chat_room_members WHERE room_id = $1 AND user_id = $2 AND left_at IS NULL)`
	var exists bool
	err := r.db.QueryRow(ctx, query, roomID, userID).Scan(&exists)
	return exists, err
}

// Message methods

func (r *repository) CreateMessage(ctx context.Context, message *Message) error {
	query := `
		INSERT INTO chat_messages (room_id, sender_id, content, reply_to_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		message.RoomID, message.SenderID, message.Content, message.ReplyToID,
	).Scan(&message.ID, &message.CreatedAt)
}

func (r *repository) GetMessageByID(ctx context.Context, id uuid.UUID) (*MessageWithDetails, error) {
	query := `
		SELECT 
			m.id, m.room_id, m.sender_id, m.content, m.reply_to_id, m.edited_at, m.deleted_at, m.created_at,
			COALESCE(s.first_name, t.first_name, mg.first_name) AS sender_first_name,
			COALESCE(s.last_name, t.last_name, mg.last_name) AS sender_last_name,
			u.email AS sender_email,
			rm.content AS reply_to_content
		FROM chat_messages m
		JOIN users u ON m.sender_id = u.id
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN teachers t ON u.id = t.user_id
		LEFT JOIN managers mg ON u.id = mg.user_id
		LEFT JOIN chat_messages rm ON m.reply_to_id = rm.id
		WHERE m.id = $1`

	var msg MessageWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&msg.ID, &msg.RoomID, &msg.SenderID, &msg.Content, &msg.ReplyToID, &msg.EditedAt, &msg.DeletedAt, &msg.CreatedAt,
		&msg.SenderFirstName, &msg.SenderLastName, &msg.SenderEmail, &msg.ReplyToContent,
	)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *repository) UpdateMessage(ctx context.Context, message *Message) error {
	query := `UPDATE chat_messages SET content = $1, edited_at = now() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, message.Content, message.ID)
	return err
}

func (r *repository) SoftDeleteMessage(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE chat_messages SET deleted_at = now(), content = '[deleted]' WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListRoomMessages(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]MessageWithDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM chat_messages WHERE room_id = $1 AND deleted_at IS NULL`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, roomID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			m.id, m.room_id, m.sender_id, m.content, m.reply_to_id, m.edited_at, m.deleted_at, m.created_at,
			COALESCE(s.first_name, t.first_name, mg.first_name) AS sender_first_name,
			COALESCE(s.last_name, t.last_name, mg.last_name) AS sender_last_name,
			u.email AS sender_email,
			rm.content AS reply_to_content
		FROM chat_messages m
		JOIN users u ON m.sender_id = u.id
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN teachers t ON u.id = t.user_id
		LEFT JOIN managers mg ON u.id = mg.user_id
		LEFT JOIN chat_messages rm ON m.reply_to_id = rm.id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []MessageWithDetails
	for rows.Next() {
		var msg MessageWithDetails
		if err := rows.Scan(
			&msg.ID, &msg.RoomID, &msg.SenderID, &msg.Content, &msg.ReplyToID, &msg.EditedAt, &msg.DeletedAt, &msg.CreatedAt,
			&msg.SenderFirstName, &msg.SenderLastName, &msg.SenderEmail, &msg.ReplyToContent,
		); err != nil {
			return nil, 0, err
		}
		messages = append(messages, msg)
	}
	return messages, total, nil
}
