package chat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ap1-final-mini-moodle/tests/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestWebSocketHandler_HandleWebSocket_Unauthorized(t *testing.T) {
	// Setup
	mockRepo := &mockRepository{}
	svc := NewService(mockRepo)

	// We don't actually need a real hub for this test since we're testing auth failure
	// before the upgrade happens.
	wsHandler := NewWebSocketHandler(svc, nil)

	router := gin.New()
	// No auth middleware - simulating unauthenticated request
	router.GET("/ws/chat/:roomId", wsHandler.HandleWebSocket)

	tests := []struct {
		name           string
		roomID         string
		setupContext   func(*gin.Context)
		wantStatusCode int
	}{
		{
			name:           "no auth context returns 401",
			roomID:         uuid.New().String(),
			setupContext:   nil, // No auth setup
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ws/chat/"+tt.roomID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d. body: %s", w.Code, tt.wantStatusCode, w.Body.String())
			}
		})
	}
}

func TestWebSocketHandler_HandleWebSocket_BadRoomID(t *testing.T) {
	// Setup
	mockRepo := &mockRepository{}
	svc := NewService(mockRepo)
	wsHandler := NewWebSocketHandler(svc, nil)

	router := gin.New()
	// Add auth middleware to set user context
	router.Use(testutil.SetAuthContext(testutil.TestUserID1, "test@test.com", "student"))
	router.GET("/ws/chat/:roomId", wsHandler.HandleWebSocket)

	tests := []struct {
		name           string
		roomID         string
		wantStatusCode int
	}{
		{
			name:           "invalid room ID format",
			roomID:         "not-a-uuid",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "empty room ID",
			roomID:         "",
			wantStatusCode: http.StatusNotFound, // Gin returns 404 for missing param
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/ws/chat/" + tt.roomID
			if tt.roomID == "" {
				path = "/ws/chat/"
			}
			req := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d. body: %s", w.Code, tt.wantStatusCode, w.Body.String())
			}
		})
	}
}

func TestWebSocketHandler_HandleWebSocket_RoomNotFound(t *testing.T) {
	roomID := uuid.New()

	// Setup - room not found (returns pgx.ErrNoRows)
	mockRepo := &mockRepository{
		// getRoomByIDFn is not set, so it will return nil, pgx.ErrNoRows by default
	}
	svc := NewService(mockRepo)
	wsHandler := NewWebSocketHandler(svc, nil)

	router := gin.New()
	router.Use(testutil.SetAuthContext(testutil.TestUserID1, "test@test.com", "student"))
	router.GET("/ws/chat/:roomId", wsHandler.HandleWebSocket)

	req := httptest.NewRequest(http.MethodGet, "/ws/chat/"+roomID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expecting 404 Not Found since room doesn't exist
	if w.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d. body: %s", w.Code, http.StatusNotFound, w.Body.String())
	}
}

func TestWebSocketHandler_HandleWebSocket_Forbidden(t *testing.T) {
	roomID := uuid.New()

	// Setup - room exists but user is not a member
	mockRepo := &mockRepository{
		getRoomByIDFn: func(ctx context.Context, id uuid.UUID) (*RoomWithDetails, error) {
			return &RoomWithDetails{
				Room: Room{ID: roomID, Type: RoomTypeGroup},
			}, nil
		},
		isMemberFn: func(ctx context.Context, rID, uID uuid.UUID) (bool, error) {
			return false, nil // Not a member
		},
	}
	svc := NewService(mockRepo)
	wsHandler := NewWebSocketHandler(svc, nil)

	router := gin.New()
	router.Use(testutil.SetAuthContext(testutil.TestUserID1, "test@test.com", "student"))
	router.GET("/ws/chat/:roomId", wsHandler.HandleWebSocket)

	req := httptest.NewRequest(http.MethodGet, "/ws/chat/"+roomID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expecting 403 Forbidden
	if w.Code != http.StatusForbidden {
		t.Errorf("got status %d, want %d. body: %s", w.Code, http.StatusForbidden, w.Body.String())
	}
}
