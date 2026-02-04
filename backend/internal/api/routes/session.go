// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/session"
	"github.com/gin-gonic/gin"
)

// SessionRouter handles session-related routes
type SessionRouter struct {
	sessionHandler *session.Handler
}

// NewSessionRouter creates a new session router
func NewSessionRouter(sessionHandler *session.Handler) *SessionRouter {
	return &SessionRouter{
		sessionHandler: sessionHandler,
	}
}

// SetupRoutes configures session routes
func (sr *SessionRouter) SetupRoutes(api *gin.RouterGroup) {
	sessions := api.Group("/sessions")
	{
		sessions.GET("", sr.sessionHandler.ListMySessions)
		sessions.DELETE("/:id", sr.sessionHandler.RevokeSession)
		sessions.POST("/revoke-all", sr.sessionHandler.RevokeAllSessions)
	}
}
