package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ap1-final-mini-moodle/internal/domain/chat"
	"github.com/ap1-final-mini-moodle/internal/shared/websocket"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	chatUC "github.com/ap1-final-mini-moodle/internal/usecase/chat"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

type ChatHandler struct {
	service     *chatUC.Service
	authService *authUC.Service
	hub         *websocket.Hub
}

func NewChatHandler(service *chatUC.Service, authService *authUC.Service, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{service: service, authService: authService, hub: hub}
}

// ─── rooms ──────────────────────────────────────────────────

func (h *ChatHandler) ListRooms(c *gin.Context) {
	userID := c.GetString("userID")
	rooms, err := h.service.ListRooms(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rooms == nil {
		rooms = []*chat.Room{}
	}
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
		"total": len(rooms),
		"page":  1,
		"limit": len(rooms),
		"total_pages": 1,
	})
}

func (h *ChatHandler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	room, err := h.service.GetRoom(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

type CreateRoomRequest struct {
	Name           *string  `json:"name"`
	Type           string   `json:"type" binding:"required"`
	CourseID       *string  `json:"course_id"`
	Members        []string `json:"members"`
	ParticipantIDs []string `json:"participant_ids"`
}

func (h *ChatHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("userID")

	// Combine members and participant_ids
	participantIDs := req.ParticipantIDs
	if len(participantIDs) == 0 {
		participantIDs = req.Members
	}

	room, err := h.service.CreateRoom(&chat.CreateRoomInput{
		RoomType:       req.Type,
		CourseID:       req.CourseID,
		Name:           req.Name,
		ParticipantIDs: participantIDs,
		CreatorUserID:  userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (h *ChatHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteRoom(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

type DirectMessageRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *ChatHandler) GetOrCreateDM(c *gin.Context) {
	var req DirectMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	currentUserID := c.GetString("userID")

	room, err := h.service.GetOrCreateDM(currentUserID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

// ─── messages ───────────────────────────────────────────────

func (h *ChatHandler) ListMessages(c *gin.Context) {
	roomID := c.Param("id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	messages, total, err := h.service.ListMessages(roomID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if messages == nil {
		messages = []*chat.Message{}
	}

	totalPages := (total + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"messages":    messages,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

type SendMessageRequest struct {
	Content    string  `json:"content" binding:"required"`
	ReplyToID  *string `json:"reply_to_id"`
	FileURL    *string `json:"file_url"`
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	roomID := c.Param("id")
	userID := c.GetString("userID")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.service.SendMessage(&chat.CreateMessageInput{
		RoomID:       roomID,
		SenderUserID: userID,
		Content:      req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Broadcast via WebSocket to OTHER users in the room (not the sender)
	roomUUID, parseErr := uuid.Parse(roomID)
	if parseErr == nil {
		senderUUID, _ := uuid.Parse(userID)
		outMsg := &websocket.OutboundMessage{
			Type:            "message",
			ID:              msg.ID,
			RoomID:          msg.RoomID,
			SenderUserID:    msg.SenderUserID,
			Content:         msg.Content,
			CreatedAt:       msg.CreatedAt,
			SenderFirstName: msg.SenderFirstName,
			SenderLastName:  msg.SenderLastName,
			SenderEmail:     msg.SenderEmail,
		}
		h.hub.BroadcastExclude(roomUUID, outMsg, senderUUID)
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	messageID := c.Param("messageId")
	if err := h.service.DeleteMessage(messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ─── members ────────────────────────────────────────────────

func (h *ChatHandler) ListMembers(c *gin.Context) {
	roomID := c.Param("id")
	members, err := h.service.ListMembers(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if members == nil {
		members = []*chat.Member{}
	}
	c.JSON(http.StatusOK, gin.H{"members": members})
}

type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *ChatHandler) AddMember(c *gin.Context) {
	roomID := c.Param("id")
	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AddMember(roomID, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"room_id": roomID, "user_id": req.UserID})
}

func (h *ChatHandler) RemoveMember(c *gin.Context) {
	roomID := c.Param("id")
	userID := c.Param("userId")
	if err := h.service.RemoveMember(roomID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *ChatHandler) LeaveRoom(c *gin.Context) {
	roomID := c.Param("id")
	userID := c.GetString("userID")
	if err := h.service.RemoveMember(roomID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "left room"})
}

// ─── users search ───────────────────────────────────────────

func (h *ChatHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	users, err := h.service.SearchUsers(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if users == nil {
		users = []chat.Participant{}
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// ─── WebSocket ──────────────────────────────────────────────

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in dev
	},
}

func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	roomID := c.Param("roomId")

	// Auth via query param token
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := h.authService.VerifyAccessToken(token)
	if err != nil {
		log.Printf("[WS] Token verification failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	log.Printf("[WS] Token verified for user: %s", claims.UserID)

	userIDParsed, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	roomIDParsed, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	client := websocket.NewClient(userIDParsed, roomIDParsed, conn, h.hub)
	h.hub.Register(client)

	// Read/write pumps
	go client.WritePump()
	go client.ReadPump(func(cl *websocket.Client, inMsg *websocket.InboundMessage) {
		switch inMsg.Type {
		case "message":
			// Save message to DB then broadcast
			msg, err := h.service.SendMessage(&chat.CreateMessageInput{
				RoomID:       roomID,
				SenderUserID: claims.UserID,
				Content:      inMsg.Content,
			})
			if err != nil {
				cl.SendError("failed to send message")
				return
			}
			outMsg := websocket.NewChatMessage(
				msg.ID, msg.RoomID, msg.SenderUserID, msg.Content, msg.CreatedAt,
			)
			outMsg.SenderFirstName = msg.SenderFirstName
			outMsg.SenderLastName = msg.SenderLastName
			outMsg.SenderEmail = msg.SenderEmail
			h.hub.Broadcast(roomIDParsed, outMsg)

		case "typing":
			// Broadcast typing indicator to others in room
			h.hub.Broadcast(roomIDParsed, &websocket.OutboundMessage{
				Type:         "typing",
				RoomID:       roomID,
				SenderUserID: claims.UserID,
			})
		}
	})
}
