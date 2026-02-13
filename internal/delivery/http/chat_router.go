package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type ChatModule struct {
	handler     *ChatHandler
	authService *authUC.Service
}

func NewChatModule(handler *ChatHandler, authService *authUC.Service) *ChatModule {
	return &ChatModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *ChatModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	// WebSocket endpoint (auth via query param token)
	api.GET("/ws/chat/:roomId", m.handler.HandleWebSocket)

	// REST endpoints (auth via Bearer token)
	chat := api.Group("/chat")
	chat.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// User search (for new message dialog)
		chat.GET("/users/search", m.handler.SearchUsers)

		// Rooms
		chat.GET("/rooms", m.handler.ListRooms)
		chat.POST("/rooms", m.handler.CreateRoom)
		chat.POST("/rooms/direct", m.handler.GetOrCreateDM)
		chat.GET("/rooms/:id", m.handler.GetRoom)
		chat.DELETE("/rooms/:id", m.handler.DeleteRoom)

		// Messages
		chat.GET("/rooms/:id/messages", m.handler.ListMessages)
		chat.POST("/rooms/:id/messages", m.handler.SendMessage)
		chat.DELETE("/rooms/:id/messages/:messageId", m.handler.DeleteMessage)

		// Members
		chat.GET("/rooms/:id/members", m.handler.ListMembers)
		chat.POST("/rooms/:id/members", m.handler.AddMember)
		chat.DELETE("/rooms/:id/members/:userId", m.handler.RemoveMember)
		chat.POST("/rooms/:id/leave", m.handler.LeaveRoom)
	}
}
