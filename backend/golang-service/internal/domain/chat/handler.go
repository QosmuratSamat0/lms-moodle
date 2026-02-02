package chat

import (
	"net/http"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles chat-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new chat handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreateRoom creates a new chat room
// @Summary Create chat room
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRoomRequest true "Room details"
// @Success 201 {object} RoomResponse
// @Router /chat/rooms [post]
func (h *Handler) CreateRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	room, err := h.service.CreateRoom(c.Request.Context(), userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, room)
}

// GetRoom gets a chat room
// @Summary Get chat room
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} RoomResponse
// @Router /chat/rooms/{id} [get]
func (h *Handler) GetRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	room, err := h.service.GetRoom(c.Request.Context(), roomID, userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, room)
}

// UpdateRoom updates a chat room
// @Summary Update chat room
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Param id path string true "Room ID"
// @Param request body UpdateRoomRequest true "Update details"
// @Success 200 {object} utils.Response
// @Router /chat/rooms/{id} [put]
func (h *Handler) UpdateRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.UpdateRoom(c.Request.Context(), roomID, userID, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "room updated successfully"})
}

// DeleteRoom deletes a chat room
// @Summary Delete chat room
// @Tags chat
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 204
// @Router /chat/rooms/{id} [delete]
func (h *Handler) DeleteRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	if err := h.service.DeleteRoom(c.Request.Context(), roomID, userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListMyRooms lists user's chat rooms
// @Summary List my chat rooms
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} RoomListResponse
// @Router /chat/rooms [get]
func (h *Handler) ListMyRooms(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListMyRooms(c.Request.Context(), userID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// DirectRoomRequest represents request to get/create a direct room
type DirectRoomRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// GetOrCreateDirectRoom gets or creates a direct room with another user
// @Summary Get or create direct chat
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body DirectRoomRequest true "Target user ID"
// @Success 200 {object} RoomResponse
// @Router /chat/rooms/direct [post]
func (h *Handler) GetOrCreateDirectRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req DirectRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if req.UserID == uuid.Nil {
		utils.BadRequest(c, "user_id is required")
		return
	}

	if req.UserID == userID {
		utils.BadRequest(c, "cannot create direct chat with yourself")
		return
	}

	room, err := h.service.GetOrCreateDirectRoom(c.Request.Context(), userID, req.UserID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, room)
}

// AddMember adds a member to a room
// @Summary Add member to room
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Param id path string true "Room ID"
// @Param request body AddMemberRequest true "Member details"
// @Success 200 {object} utils.Response
// @Router /chat/rooms/{id}/members [post]
func (h *Handler) AddMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.AddMember(c.Request.Context(), roomID, userID, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "member added successfully"})
}

// RemoveMember removes a member from a room
// @Summary Remove member from room
// @Tags chat
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Param user_id path string true "User ID"
// @Success 204
// @Router /chat/rooms/{id}/members/{user_id} [delete]
func (h *Handler) RemoveMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	targetIDStr := c.Param("user_id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid user ID")
		return
	}

	if err := h.service.RemoveMember(c.Request.Context(), roomID, userID, targetID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// LeaveRoom leaves a chat room
// @Summary Leave chat room
// @Tags chat
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 204
// @Router /chat/rooms/{id}/leave [post]
func (h *Handler) LeaveRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	if err := h.service.LeaveRoom(c.Request.Context(), roomID, userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListMembers lists room members
// @Summary List room members
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {array} MemberResponse
// @Router /chat/rooms/{id}/members [get]
func (h *Handler) ListMembers(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	members, err := h.service.ListMembers(c.Request.Context(), roomID, userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, members)
}

// SendMessage sends a message
// @Summary Send message
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Param request body SendMessageRequest true "Message content"
// @Success 201 {object} MessageResponse
// @Router /chat/rooms/{id}/messages [post]
func (h *Handler) SendMessage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	message, err := h.service.SendMessage(c.Request.Context(), roomID, userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, message)
}

// UpdateMessage updates a message
// @Summary Update message
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Param message_id path string true "Message ID"
// @Param request body UpdateMessageRequest true "Message content"
// @Success 200 {object} utils.Response
// @Router /chat/messages/{message_id} [put]
func (h *Handler) UpdateMessage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	messageIDStr := c.Param("message_id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid message ID")
		return
	}

	var req UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.UpdateMessage(c.Request.Context(), messageID, userID, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "message updated successfully"})
}

// DeleteMessage deletes a message
// @Summary Delete message
// @Tags chat
// @Security BearerAuth
// @Param message_id path string true "Message ID"
// @Success 204
// @Router /chat/messages/{message_id} [delete]
func (h *Handler) DeleteMessage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	messageIDStr := c.Param("message_id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid message ID")
		return
	}

	if err := h.service.DeleteMessage(c.Request.Context(), messageID, userID); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListMessages lists room messages
// @Summary List messages
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Param id path string true "Room ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} MessageListResponse
// @Router /chat/rooms/{id}/messages [get]
func (h *Handler) ListMessages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid room ID")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListMessages(c.Request.Context(), roomID, userID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}
