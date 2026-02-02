// Package routes contains route registration helpers
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/chat"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/websocket"
	"github.com/gin-gonic/gin"
)

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
