package chat

import (
	"context"
	"net/http"
	"strings"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

const (
	maxContentLength = 5000
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, check specific origins
		return true
	},
}

// WebSocketHandler handles WebSocket connections for chat.
type WebSocketHandler struct {
	service Service
	hub     *websocket.Hub
}

// NewWebSocketHandler creates a new WebSocket handler.
func NewWebSocketHandler(service Service, hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		service: service,
		hub:     hub,
	}
}

// HandleWebSocket handles the WebSocket upgrade and connection.
// GET /ws/chat/:roomId
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Extract user from auth context
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	userRole := middleware.GetUserRole(c)

	// Parse room ID
	roomIDStr := c.Param("roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	// Authorize room access
	ctx := c.Request.Context()
	if err := h.service.AuthorizeRoomAccess(ctx, roomID, userID, userRole); err != nil {
		handleWSAuthError(c, err)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Upgrader already writes error response
		return
	}

	// Create client and register with hub
	client := websocket.NewClient(userID, roomID, conn, h.hub)
	h.hub.Register(client)

	// Start pumps
	go client.WritePump()
	go client.ReadPump(func(cl *websocket.Client, msg *websocket.InboundMessage) {
		h.handleInboundMessage(ctx, cl, msg)
	})
}

// handleInboundMessage processes incoming WebSocket messages.
func (h *WebSocketHandler) handleInboundMessage(ctx context.Context, client *websocket.Client, msg *websocket.InboundMessage) {
	switch msg.Type {
	case "message":
		h.handleChatMessage(ctx, client, msg)
	default:
		client.SendError("unknown message type")
	}
}

// handleChatMessage processes chat messages.
func (h *WebSocketHandler) handleChatMessage(ctx context.Context, client *websocket.Client, msg *websocket.InboundMessage) {
	// Validate content
	content := strings.TrimSpace(msg.Content)
	if content == "" {
		client.SendError("message content cannot be empty")
		return
	}
	if len(content) > maxContentLength {
		client.SendError("message content too long")
		return
	}

	// Save and broadcast message
	outMsg, err := h.service.SaveAndBroadcastMessage(ctx, client.RoomID, client.UserID, content, h.hub)
	if err != nil {
		client.SendError("failed to send message")
		return
	}

	// Broadcast is handled by service
	_ = outMsg
}

// handleWSAuthError handles authorization errors for WebSocket upgrade.
func handleWSAuthError(c *gin.Context, err error) {
	// Before upgrade, we can still send HTTP errors
	if strings.Contains(err.Error(), "not found") {
		utils.Error(c, http.StatusNotFound, "NOT_FOUND", "chat room not found")
		return
	}
	utils.Error(c, http.StatusForbidden, "FORBIDDEN", err.Error())
}
