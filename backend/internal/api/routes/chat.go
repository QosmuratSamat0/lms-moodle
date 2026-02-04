// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/chat"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/websocket"
	"github.com/gin-gonic/gin"
)

// ChatRouter handles chat-related routes
type ChatRouter struct {
	chatHandler *chat.Handler
}

// NewChatRouter creates a new chat router
func NewChatRouter(chatHandler *chat.Handler) *ChatRouter {
	return &ChatRouter{
		chatHandler: chatHandler,
	}
}

// SetupRoutes configures chat routes
func (cr *ChatRouter) SetupRoutes(api *gin.RouterGroup) {
	chatRoutes := api.Group("/chat")
	{
		chatRoutes.GET("/rooms", cr.chatHandler.ListMyRooms)
		chatRoutes.GET("/rooms/:id", cr.chatHandler.GetRoom)
		chatRoutes.POST("/rooms", cr.chatHandler.CreateRoom)
		chatRoutes.PUT("/rooms/:id", cr.chatHandler.UpdateRoom)
		chatRoutes.DELETE("/rooms/:id", cr.chatHandler.DeleteRoom)
		chatRoutes.POST("/rooms/direct", cr.chatHandler.GetOrCreateDirectRoom)
		chatRoutes.GET("/rooms/:id/messages", cr.chatHandler.ListMessages)
		chatRoutes.POST("/rooms/:id/messages", cr.chatHandler.SendMessage)
		chatRoutes.PUT("/messages/:id", cr.chatHandler.UpdateMessage)
		chatRoutes.DELETE("/messages/:id", cr.chatHandler.DeleteMessage)
		chatRoutes.POST("/rooms/:id/members", cr.chatHandler.AddMember)
		chatRoutes.DELETE("/rooms/:id/members/:user_id", cr.chatHandler.RemoveMember)
		chatRoutes.POST("/rooms/:id/leave", cr.chatHandler.LeaveRoom)
		chatRoutes.GET("/rooms/:id/members", cr.chatHandler.ListMembers)
	}
}

// SetupChatWebSocket sets up the WebSocket route for chat.
// This should be called with the engine directly (not a router group)
// since WebSocket routes often need different handling.
func SetupChatWebSocket(engine *gin.Engine, auth *middleware.AuthMiddleware, service chat.Service, hub *websocket.Hub) {
	wsHandler := chat.NewWebSocketHandler(service, hub)

	// WebSocket route with auth
	ws := engine.Group("/api/v1/ws")
	ws.Use(auth.Authenticate())
	{
		ws.GET("/chat/:roomId", wsHandler.HandleWebSocket)
	}
}
