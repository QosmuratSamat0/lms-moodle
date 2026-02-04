package schedule

import (
	"net/http"
	"time"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles schedule-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new schedule handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create creates a new schedule event
// @Summary Create schedule event
// @Tags schedule
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateEventRequest true "Event details"
// @Success 201 {object} EventResponse
// @Router /schedule/events [post]
func (h *Handler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	event, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, event)
}

// GetByID gets a schedule event by ID
// @Summary Get schedule event
// @Tags schedule
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} EventResponse
// @Router /schedule/events/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid event ID")
		return
	}

	event, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, event)
}

// Update updates a schedule event
// @Summary Update schedule event
// @Tags schedule
// @Security BearerAuth
// @Accept json
// @Param id path string true "Event ID"
// @Param request body UpdateEventRequest true "Update details"
// @Success 200 {object} utils.Response
// @Router /schedule/events/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid event ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(c.Request.Context(), id, userID, role, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "event updated successfully"})
}

// Delete deletes a schedule event
// @Summary Delete schedule event
// @Tags schedule
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Success 204
// @Router /schedule/events/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid event ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	if err := h.service.Delete(c.Request.Context(), id, userID, role); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListByCourse lists events for a course
// @Summary List course schedule
// @Tags schedule
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Param event_type query string false "Event type filter"
// @Param start_date query string false "Start date filter (RFC3339)"
// @Param end_date query string false "End date filter (RFC3339)"
// @Success 200 {object} EventListResponse
// @Router /courses/{id}/schedule [get]
func (h *Handler) ListByCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	pagination := utils.GetPaginationFromContext(c)
	filter := h.parseFilter(c)

	list, err := h.service.ListByCourse(c.Request.Context(), courseID, filter, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// ListMySchedule lists user's schedule
// @Summary List my schedule
// @Tags schedule
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Param event_type query string false "Event type filter"
// @Param start_date query string false "Start date filter (RFC3339)"
// @Param end_date query string false "End date filter (RFC3339)"
// @Success 200 {object} EventListResponse
// @Router /schedule [get]
func (h *Handler) ListMySchedule(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)
	filter := h.parseFilter(c)

	list, err := h.service.ListMySchedule(c.Request.Context(), userID, filter, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// GetCalendar gets calendar events for a date range
// @Summary Get calendar events
// @Tags schedule
// @Security BearerAuth
// @Produce json
// @Param start query string true "Start date (RFC3339)"
// @Param end query string true "End date (RFC3339)"
// @Success 200 {array} EventResponse
// @Router /schedule/calendar [get]
func (h *Handler) GetCalendar(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		utils.BadRequest(c, "start and end dates are required")
		return
	}

	startDate, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		utils.BadRequest(c, "invalid start date format")
		return
	}

	endDate, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		utils.BadRequest(c, "invalid end date format")
		return
	}

	events, err := h.service.GetCalendarEvents(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, events)
}

func (h *Handler) parseFilter(c *gin.Context) *EventFilter {
	filter := &EventFilter{}

	if eventType := c.Query("event_type"); eventType != "" {
		filter.EventType = &eventType
	}
	if startDate := c.Query("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	return filter
}
