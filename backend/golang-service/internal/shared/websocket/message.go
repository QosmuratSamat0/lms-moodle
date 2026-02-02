package websocket

import (
	"encoding/json"
	"time"
)

// InboundMessage represents a message received from client.
type InboundMessage struct {
	Type    string          `json:"type"`
	Content string          `json:"content,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// OutboundMessage represents a message sent to clients.
type OutboundMessage struct {
	Type         string    `json:"type"`
	ID           string    `json:"id,omitempty"`
	RoomID       string    `json:"roomId,omitempty"`
	SenderUserID string    `json:"senderUserId,omitempty"`
	Content      string    `json:"content,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	Error        string    `json:"error,omitempty"`
}

// NewErrorMessage creates an error outbound message.
func NewErrorMessage(errMsg string) *OutboundMessage {
	return &OutboundMessage{
		Type:  "error",
		Error: errMsg,
	}
}

// NewChatMessage creates a chat outbound message.
func NewChatMessage(id, roomID, senderUserID, content string, createdAt time.Time) *OutboundMessage {
	return &OutboundMessage{
		Type:         "message",
		ID:           id,
		RoomID:       roomID,
		SenderUserID: senderUserID,
		Content:      content,
		CreatedAt:    createdAt,
	}
}
