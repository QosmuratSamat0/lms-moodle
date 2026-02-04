package session

import (
	"net/http"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles session-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new session handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ListMySessions lists current user's active sessions
// @Summary List my sessions
// @Tags sessions
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} SessionListResponse
// @Router /sessions [get]
func (h *Handler) ListMySessions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListMySessions(c.Request.Context(), userID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// RevokeSession revokes a specific session
// @Summary Revoke session
// @Tags sessions
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 204
// @Router /sessions/{id} [delete]
func (h *Handler) RevokeSession(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid session ID")
		return
	}

	if err := h.service.Revoke(c.Request.Context(), id, userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// RevokeAllSessions revokes all sessions for current user
// @Summary Revoke all sessions
// @Tags sessions
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /sessions/revoke-all [post]
func (h *Handler) RevokeAllSessions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	if err := h.service.RevokeAll(c.Request.Context(), userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "all sessions revoked"})
}
