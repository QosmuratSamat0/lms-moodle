//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/chat"
	"github.com/MaqsattoTeam/aLMS/golang-service/tests/testutil"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocket message types for testing.
type wsInboundMessage struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
}

type wsOutboundMessage struct {
	Type         string `json:"type"`
	ID           string `json:"id,omitempty"`
	RoomID       string `json:"roomId,omitempty"`
	SenderUserID string `json:"senderUserId,omitempty"`
	Content      string `json:"content,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	Error        string `json:"error,omitempty"`
}

func TestWebSocketChat_E2E(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)

	// Setup test data
	testutil.TruncateAll(t, pool)

	// Register and login a user to get auth token
	client := &http.Client{Timeout: 10 * time.Second}

	// Register user
	registerReq := map[string]interface{}{
		"email":    "wsuser@test.com",
		"password": "testpassword123",
		"role":     "student",
	}
	bodyBytes, _ := json.Marshal(registerReq)
	resp, err := client.Post(
		getBaseURL()+"/api/v1/auth/register",
		"application/json",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	resp.Body.Close()

	// Login to get token
	loginReq := map[string]interface{}{
		"email":    "wsuser@test.com",
		"password": "testpassword123",
	}
	bodyBytes, _ = json.Marshal(loginReq)
	resp, err = client.Post(
		getBaseURL()+"/api/v1/auth/login",
		"application/json",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
			User        struct {
				ID string `json:"id"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}

	accessToken := loginResp.Data.AccessToken
	userID, _ := uuid.Parse(loginResp.Data.User.ID)

	if accessToken == "" {
		t.Fatal("no access token received from login")
	}

	// Create chat room via API or directly in DB
	chatRepo := chat.NewRepository(pool)
	room := &chat.Room{
		Name:      stringPtr("E2E Test Room"),
		Type:      chat.RoomTypeGroup,
		CreatedBy: userID,
	}
	roomCtx := context.Background()
	if err := chatRepo.CreateRoom(roomCtx, room); err != nil {
		t.Fatalf("failed to create chat room: %v", err)
	}

	// Add user as member
	member := &chat.Member{
		RoomID: room.ID,
		UserID: userID,
		Role:   "admin",
	}
	if err := chatRepo.AddMember(roomCtx, member); err != nil {
		t.Fatalf("failed to add member: %v", err)
	}

	// Also use fixtures to ensure student profile exists
	fixtures.CreateStudent(roomCtx, t, userID, "WS", "User", "CS-E2E")

	t.Run("connect and send message", func(t *testing.T) {
		// Build WebSocket URL
		wsURL := strings.Replace(getBaseURL(), "http://", "ws://", 1) + "/api/v1/ws/chat/" + room.ID.String()
		u, err := url.Parse(wsURL)
		if err != nil {
			t.Fatalf("failed to parse websocket URL: %v", err)
		}

		// Connect with auth header
		header := http.Header{}
		header.Add("Authorization", "Bearer "+accessToken)

		conn, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			if resp != nil {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("websocket dial failed: %v, response: %s", err, string(body))
			}
			t.Fatalf("websocket dial failed: %v", err)
		}
		defer conn.Close()

		// Send a message
		inMsg := wsInboundMessage{
			Type:    "message",
			Content: "Hello from E2E test!",
		}
		if err := conn.WriteJSON(inMsg); err != nil {
			t.Fatalf("failed to send message: %v", err)
		}

		// Read the broadcast response
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		var outMsg wsOutboundMessage
		if err := conn.ReadJSON(&outMsg); err != nil {
			t.Fatalf("failed to read response: %v", err)
		}

		// Verify response
		if outMsg.Type != "message" {
			t.Errorf("expected type 'message', got %q", outMsg.Type)
		}
		if outMsg.Content != "Hello from E2E test!" {
			t.Errorf("expected content 'Hello from E2E test!', got %q", outMsg.Content)
		}
		if outMsg.RoomID != room.ID.String() {
			t.Errorf("expected roomId %q, got %q", room.ID.String(), outMsg.RoomID)
		}
		if outMsg.SenderUserID != userID.String() {
			t.Errorf("expected senderUserId %q, got %q", userID.String(), outMsg.SenderUserID)
		}
		if outMsg.ID == "" {
			t.Error("expected message ID to be set")
		}
		if outMsg.CreatedAt == "" {
			t.Error("expected createdAt to be set")
		}
	})

	t.Run("unauthorized connection rejected", func(t *testing.T) {
		// Build WebSocket URL
		wsURL := strings.Replace(getBaseURL(), "http://", "ws://", 1) + "/api/v1/ws/chat/" + room.ID.String()
		u, err := url.Parse(wsURL)
		if err != nil {
			t.Fatalf("failed to parse websocket URL: %v", err)
		}

		// Connect without auth header
		conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err == nil {
			conn.Close()
			t.Fatal("expected connection to be rejected without auth")
		}

		if resp != nil && resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("non-member connection rejected", func(t *testing.T) {
		// Register another user
		registerReq := map[string]interface{}{
			"email":    "wsuser2@test.com",
			"password": "testpassword123",
			"role":     "student",
		}
		bodyBytes, _ := json.Marshal(registerReq)
		resp, err := client.Post(
			getBaseURL()+"/api/v1/auth/register",
			"application/json",
			bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			t.Fatalf("failed to register second user: %v", err)
		}
		resp.Body.Close()

		// Login second user
		loginReq := map[string]interface{}{
			"email":    "wsuser2@test.com",
			"password": "testpassword123",
		}
		bodyBytes, _ = json.Marshal(loginReq)
		resp, err = client.Post(
			getBaseURL()+"/api/v1/auth/login",
			"application/json",
			bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			t.Fatalf("failed to login second user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("login failed for second user: %s", string(body))
			return // Skip if registration/login failed
		}

		var loginResp2 struct {
			Data struct {
				AccessToken string `json:"access_token"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&loginResp2); err != nil {
			t.Logf("failed to decode login response: %v", err)
			return
		}

		accessToken2 := loginResp2.Data.AccessToken
		if accessToken2 == "" {
			t.Log("no access token for second user, skipping")
			return
		}

		// Try to connect to room (user2 is not a member)
		wsURL := strings.Replace(getBaseURL(), "http://", "ws://", 1) + "/api/v1/ws/chat/" + room.ID.String()
		u, _ := url.Parse(wsURL)

		header := http.Header{}
		header.Add("Authorization", "Bearer "+accessToken2)

		conn, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err == nil {
			conn.Close()
			t.Fatal("expected connection to be rejected for non-member")
		}

		if resp != nil && resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid room ID rejected", func(t *testing.T) {
		wsURL := strings.Replace(getBaseURL(), "http://", "ws://", 1) + "/api/v1/ws/chat/not-a-uuid"
		u, _ := url.Parse(wsURL)

		header := http.Header{}
		header.Add("Authorization", "Bearer "+accessToken)

		conn, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err == nil {
			conn.Close()
			t.Fatal("expected connection to be rejected for invalid room ID")
		}

		if resp != nil && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("non-existent room rejected", func(t *testing.T) {
		nonExistentRoomID := uuid.New()
		wsURL := strings.Replace(getBaseURL(), "http://", "ws://", 1) + "/api/v1/ws/chat/" + nonExistentRoomID.String()
		u, _ := url.Parse(wsURL)

		header := http.Header{}
		header.Add("Authorization", "Bearer "+accessToken)

		conn, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err == nil {
			conn.Close()
			t.Fatal("expected connection to be rejected for non-existent room")
		}

		if resp != nil && resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("empty message rejected", func(t *testing.T) {
		// Build WebSocket URL
		wsURL := strings.Replace(getBaseURL(), "http://", "ws://", 1) + "/api/v1/ws/chat/" + room.ID.String()
		u, _ := url.Parse(wsURL)

		header := http.Header{}
		header.Add("Authorization", "Bearer "+accessToken)

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			t.Fatalf("websocket dial failed: %v", err)
		}
		defer conn.Close()

		// Send empty message
		inMsg := wsInboundMessage{
			Type:    "message",
			Content: "",
		}
		if err := conn.WriteJSON(inMsg); err != nil {
			t.Fatalf("failed to send message: %v", err)
		}

		// Read the error response
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		var outMsg wsOutboundMessage
		if err := conn.ReadJSON(&outMsg); err != nil {
			t.Fatalf("failed to read response: %v", err)
		}

		// Verify error response
		if outMsg.Type != "error" {
			t.Errorf("expected type 'error', got %q", outMsg.Type)
		}
		if outMsg.Error == "" {
			t.Error("expected error message to be set")
		}
	})
}

func stringPtr(s string) *string {
	return &s
}
