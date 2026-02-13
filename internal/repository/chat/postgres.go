package chat

import (
	"context"
	"fmt"

	"github.com/ap1-final-mini-moodle/internal/domain/chat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) chat.Repository {
	return &PostgresRepository{db: db}
}

// ─── rooms ──────────────────────────────────────────────────

func (r *PostgresRepository) CreateRoom(room *chat.Room) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO chat_rooms (id, course_id, room_type, created_at)
		 VALUES ($1, $2, $3, $4)`,
		room.ID, room.CourseID, room.RoomType, room.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetRoomByID(id string, currentUserID string) (*chat.Room, error) {
	room := &chat.Room{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, course_id, room_type, created_at
		 FROM chat_rooms WHERE id = $1`, id).
		Scan(&room.ID, &room.CourseID, &room.RoomType, &room.CreatedAt)
	if err != nil {
		return nil, err
	}

	// member count
	r.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM chat_room_members WHERE room_id=$1`, id).
		Scan(&room.MemberCount)

	// participants — exclude current user for DM rooms so header shows the other person
	if currentUserID != "" {
		participants, _ := r.getRoomParticipantsExcluding(id, currentUserID)
		room.Participants = participants
	} else {
		participants, _ := r.getRoomParticipants(id)
		room.Participants = participants
	}

	// last message
	lm := &chat.LastMessage{}
	var senderFirst *string
	err2 := r.db.QueryRow(context.Background(),
		`SELECT m.id, m.content, m.sender_user_id,
		        COALESCE(s.first_name, t.first_name, mg.first_name, ''),
		        m.created_at
		 FROM chat_messages m
		 LEFT JOIN users u ON u.id = m.sender_user_id
		 LEFT JOIN students s ON s.user_id = u.id
		 LEFT JOIN teachers t ON t.user_id = u.id
		 LEFT JOIN managers mg ON mg.user_id = u.id
		 WHERE m.room_id = $1
		 ORDER BY m.created_at DESC LIMIT 1`, id).
		Scan(&lm.ID, &lm.Content, &lm.SenderID, &senderFirst, &room.LastMessageAt)
	if err2 == nil {
		if senderFirst != nil {
			lm.SenderFirstName = *senderFirst
		}
		if room.LastMessageAt != nil {
			lm.CreatedAt = room.LastMessageAt.Format("2006-01-02T15:04:05Z")
		}
		room.LastMessage = lm
	}

	return room, nil
}

func (r *PostgresRepository) ListRoomsByUser(userID string) ([]*chat.Room, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT cr.id, cr.course_id, cr.room_type, cr.created_at
		 FROM chat_rooms cr
		 INNER JOIN chat_room_members crm ON crm.room_id = cr.id
		 WHERE crm.user_id = $1
		 ORDER BY cr.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*chat.Room
	for rows.Next() {
		room := &chat.Room{}
		if err := rows.Scan(&room.ID, &room.CourseID, &room.RoomType, &room.CreatedAt); err != nil {
			return nil, err
		}

		// member count
		r.db.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM chat_room_members WHERE room_id=$1`, room.ID).
			Scan(&room.MemberCount)

		// participants (exclude current user for DMs)
		participants, _ := r.getRoomParticipantsExcluding(room.ID, userID)
		room.Participants = participants

		// last message
		lm := &chat.LastMessage{}
		var senderFirst *string
		err2 := r.db.QueryRow(context.Background(),
			`SELECT m.id, m.content, m.sender_user_id,
			        COALESCE(s.first_name, t.first_name, mg.first_name, ''),
			        m.created_at
			 FROM chat_messages m
			 LEFT JOIN users u ON u.id = m.sender_user_id
			 LEFT JOIN students s ON s.user_id = u.id
			 LEFT JOIN teachers t ON t.user_id = u.id
			 LEFT JOIN managers mg ON mg.user_id = u.id
			 WHERE m.room_id = $1
			 ORDER BY m.created_at DESC LIMIT 1`, room.ID).
			Scan(&lm.ID, &lm.Content, &lm.SenderID, &senderFirst, &room.LastMessageAt)
		if err2 == nil {
			if senderFirst != nil {
				lm.SenderFirstName = *senderFirst
			}
			if room.LastMessageAt != nil {
				lm.CreatedAt = room.LastMessageAt.Format("2006-01-02T15:04:05Z")
			}
			room.LastMessage = lm
		}

		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (r *PostgresRepository) DeleteRoom(id string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM chat_rooms WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) FindDirectRoom(userA, userB string) (*chat.Room, error) {
	var roomID string
	err := r.db.QueryRow(context.Background(),
		`SELECT cr.id
		 FROM chat_rooms cr
		 WHERE cr.room_type = 'direct'
		   AND EXISTS (SELECT 1 FROM chat_room_members WHERE room_id = cr.id AND user_id = $1)
		   AND EXISTS (SELECT 1 FROM chat_room_members WHERE room_id = cr.id AND user_id = $2)
		 LIMIT 1`, userA, userB).Scan(&roomID)
	if err != nil {
		return nil, err
	}
	return r.GetRoomByID(roomID, userA)
}

// ─── members ────────────────────────────────────────────────

func (r *PostgresRepository) AddMember(roomID, userID string) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO chat_room_members (room_id, user_id) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`, roomID, userID)
	return err
}

func (r *PostgresRepository) RemoveMember(roomID, userID string) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM chat_room_members WHERE room_id = $1 AND user_id = $2`, roomID, userID)
	return err
}

func (r *PostgresRepository) ListMembers(roomID string) ([]*chat.Member, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT crm.room_id, crm.user_id, crm.joined_at,
		        COALESCE(s.first_name, t.first_name, mg.first_name, ''),
		        COALESCE(s.last_name, t.last_name, mg.last_name, ''),
		        u.email
		 FROM chat_room_members crm
		 JOIN users u ON u.id = crm.user_id
		 LEFT JOIN students s ON s.user_id = u.id
		 LEFT JOIN teachers t ON t.user_id = u.id
		 LEFT JOIN managers mg ON mg.user_id = u.id
		 WHERE crm.room_id = $1
		 ORDER BY crm.joined_at`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*chat.Member
	for rows.Next() {
		m := &chat.Member{}
		if err := rows.Scan(&m.RoomID, &m.UserID, &m.JoinedAt,
			&m.FirstName, &m.LastName, &m.Email); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *PostgresRepository) IsMember(roomID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM chat_room_members WHERE room_id=$1 AND user_id=$2)`,
		roomID, userID).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) GetMemberCount(roomID string) (int, error) {
	var count int
	err := r.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM chat_room_members WHERE room_id=$1`, roomID).Scan(&count)
	return count, err
}

// ─── messages ───────────────────────────────────────────────

func (r *PostgresRepository) CreateMessage(msg *chat.Message) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO chat_messages (id, room_id, sender_user_id, content, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		msg.ID, msg.RoomID, msg.SenderUserID, msg.Content, msg.CreatedAt)
	return err
}

func (r *PostgresRepository) GetMessageByID(id string) (*chat.Message, error) {
	m := &chat.Message{}
	err := r.db.QueryRow(context.Background(),
		`SELECT m.id, m.room_id, m.sender_user_id, m.content, m.created_at,
		        COALESCE(s.first_name, t.first_name, mg.first_name, ''),
		        COALESCE(s.last_name, t.last_name, mg.last_name, ''),
		        u.email
		 FROM chat_messages m
		 LEFT JOIN users u ON u.id = m.sender_user_id
		 LEFT JOIN students s ON s.user_id = u.id
		 LEFT JOIN teachers t ON t.user_id = u.id
		 LEFT JOIN managers mg ON mg.user_id = u.id
		 WHERE m.id = $1`, id).
		Scan(&m.ID, &m.RoomID, &m.SenderUserID, &m.Content, &m.CreatedAt,
			&m.SenderFirstName, &m.SenderLastName, &m.SenderEmail)
	return m, err
}

func (r *PostgresRepository) ListMessages(roomID string, limit, offset int) ([]*chat.Message, int, error) {
	if limit <= 0 {
		limit = 50
	}

	var total int
	r.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM chat_messages WHERE room_id=$1`, roomID).Scan(&total)

	rows, err := r.db.Query(context.Background(),
		`SELECT m.id, m.room_id, m.sender_user_id, m.content, m.created_at,
		        COALESCE(s.first_name, t.first_name, mg.first_name, ''),
		        COALESCE(s.last_name, t.last_name, mg.last_name, ''),
		        COALESCE(u.email, '')
		 FROM chat_messages m
		 LEFT JOIN users u ON u.id = m.sender_user_id
		 LEFT JOIN students s ON s.user_id = u.id
		 LEFT JOIN teachers t ON t.user_id = u.id
		 LEFT JOIN managers mg ON mg.user_id = u.id
		 WHERE m.room_id = $1
		 ORDER BY m.created_at DESC
		 LIMIT $2 OFFSET $3`, roomID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []*chat.Message
	for rows.Next() {
		m := &chat.Message{}
		if err := rows.Scan(&m.ID, &m.RoomID, &m.SenderUserID, &m.Content, &m.CreatedAt,
			&m.SenderFirstName, &m.SenderLastName, &m.SenderEmail); err != nil {
			return nil, 0, err
		}
		messages = append(messages, m)
	}
	return messages, total, rows.Err()
}

func (r *PostgresRepository) DeleteMessage(id string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM chat_messages WHERE id = $1`, id)
	return err
}

// ─── users search ───────────────────────────────────────────

func (r *PostgresRepository) SearchUsers(query string, limit int) ([]chat.Participant, error) {
	if limit <= 0 {
		limit = 20
	}
	if query == "" {
		return []chat.Participant{}, nil
	}
	pattern := fmt.Sprintf("%%%s%%", query)

	rows, err := r.db.Query(context.Background(),
		`SELECT u.id,
		        COALESCE(s.first_name, t.first_name, m.first_name, '') AS first_name,
		        COALESCE(s.last_name, t.last_name, m.last_name, '') AS last_name,
		        u.email
		 FROM users u
		 LEFT JOIN students s ON s.user_id = u.id
		 LEFT JOIN teachers t ON t.user_id = u.id
		 LEFT JOIN managers m ON m.user_id = u.id
		 WHERE u.is_active = true
		   AND (
		     LOWER(COALESCE(s.first_name, t.first_name, m.first_name, '')) LIKE LOWER($1)
		     OR LOWER(COALESCE(s.last_name, t.last_name, m.last_name, '')) LIKE LOWER($1)
		     OR LOWER(u.email) LIKE LOWER($1)
		     OR LOWER(
		       COALESCE(s.first_name, t.first_name, m.first_name, '') || ' ' ||
		       COALESCE(s.last_name, t.last_name, m.last_name, '')
		     ) LIKE LOWER($1)
		   )
		 ORDER BY first_name, last_name
		 LIMIT $2`, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []chat.Participant
	for rows.Next() {
		u := chat.Participant{}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// ─── helpers ────────────────────────────────────────────────

func (r *PostgresRepository) getRoomParticipants(roomID string) ([]chat.Participant, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT u.id,
		        COALESCE(s.first_name, t.first_name, mg.first_name, '') AS first_name,
		        COALESCE(s.last_name, t.last_name, mg.last_name, '') AS last_name,
		        u.email
		 FROM chat_room_members crm
		 JOIN users u ON u.id = crm.user_id
		 LEFT JOIN students s ON s.user_id = u.id
		 LEFT JOIN teachers t ON t.user_id = u.id
		 LEFT JOIN managers mg ON mg.user_id = u.id
		 WHERE crm.room_id = $1`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []chat.Participant
	for rows.Next() {
		p := chat.Participant{}
		if err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}

func (r *PostgresRepository) getRoomParticipantsExcluding(roomID, excludeUserID string) ([]chat.Participant, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT u.id,
		        COALESCE(s.first_name, t.first_name, mg.first_name, '') AS first_name,
		        COALESCE(s.last_name, t.last_name, mg.last_name, '') AS last_name,
		        u.email
		 FROM chat_room_members crm
		 JOIN users u ON u.id = crm.user_id
		 LEFT JOIN students s ON s.user_id = u.id
		 LEFT JOIN teachers t ON t.user_id = u.id
		 LEFT JOIN managers mg ON mg.user_id = u.id
		 WHERE crm.room_id = $1 AND crm.user_id != $2`, roomID, excludeUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []chat.Participant
	for rows.Next() {
		p := chat.Participant{}
		if err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}
