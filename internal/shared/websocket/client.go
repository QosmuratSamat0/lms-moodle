package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 4096
)

type Client struct {
	ID     string
	UserID uuid.UUID
	RoomID uuid.UUID
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	mu     sync.Mutex
}

func NewClient(userID, roomID uuid.UUID, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:     uuid.New().String(),
		UserID: userID,
		RoomID: roomID,
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    hub,
	}
}

func (c *Client) ReadPump(onMessage func(client *Client, msg *InboundMessage)) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		var inMsg InboundMessage
		if err := json.Unmarshal(message, &inMsg); err != nil {
			c.SendError("invalid message format")
			continue
		}

		if onMessage != nil {
			onMessage(c, &inMsg)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Send(msg *OutboundMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.send <- data:
	default:
		log.Printf("client %s buffer full, dropping message", c.ID)
	}
}

func (c *Client) SendError(errMsg string) {
	c.Send(NewErrorMessage(errMsg))
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	close(c.send)
}
