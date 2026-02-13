package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	rooms map[uuid.UUID]map[*Client]bool

	register chan *Client

	unregister chan *Client

	broadcast chan *broadcastMessage

	mu sync.RWMutex
}

type broadcastMessage struct {
	roomID        uuid.UUID
	message       *OutboundMessage
	excludeUserID *uuid.UUID
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uuid.UUID]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *broadcastMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()
			log.Printf("client %s joined room %s", client.UserID, client.RoomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					client.Close()
					if len(clients) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("client %s left room %s", client.UserID, client.RoomID)

		case bm := <-h.broadcast:
			h.mu.RLock()
			clients := h.rooms[bm.roomID]
			h.mu.RUnlock()

			data, err := json.Marshal(bm.message)
			if err != nil {
				log.Printf("failed to marshal broadcast message: %v", err)
				continue
			}

			for client := range clients {
				// Skip excluded user
				if bm.excludeUserID != nil && client.UserID == *bm.excludeUserID {
					continue
				}
				select {
				case client.send <- data:
				default:
					h.mu.Lock()
					delete(h.rooms[client.RoomID], client)
					client.Close()
					h.mu.Unlock()
				}
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(roomID uuid.UUID, msg *OutboundMessage) {
	h.broadcast <- &broadcastMessage{
		roomID:  roomID,
		message: msg,
	}
}

func (h *Hub) BroadcastExclude(roomID uuid.UUID, msg *OutboundMessage, excludeUserID uuid.UUID) {
	h.broadcast <- &broadcastMessage{
		roomID:        roomID,
		message:       msg,
		excludeUserID: &excludeUserID,
	}
}

func (h *Hub) GetClientsInRoom(roomID uuid.UUID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}

func (h *Hub) IsUserInRoom(roomID, userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients := h.rooms[roomID]
	for client := range clients {
		if client.UserID == userID {
			return true
		}
	}
	return false
}
