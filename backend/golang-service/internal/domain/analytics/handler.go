package analytics

import (
	"strconv"
	"time"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles analytics-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new analytics handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// TrackEvent tracks an analytics event
// @Summary Track event
// @Tags analytics
// @Security BearerAuth
// @Accept json
// @Param request body TrackEventRequest true "Event details"
// @Success 200 {object} utils.Response
// @Router /analytics/track [post]
func (h *Handler) TrackEvent(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req TrackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	if err := h.service.TrackEvent(c.Request.Context(), userID, &req, &userAgent, &ipAddress); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "event tracked"})
}

// GetCourseStats gets course statistics
// @Summary Get course statistics
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} CourseStatResponse
// @Router /analytics/courses/{id} [get]
func (h *Handler) GetCourseStats(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	stats, err := h.service.GetCourseStats(c.Request.Context(), courseID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, stats)
}

// GetMyStats gets current user's statistics
// @Summary Get my statistics
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserStatResponse
// @Router /analytics/me [get]
func (h *Handler) GetMyStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	stats, err := h.service.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, stats)
}

// GetUserStats gets a specific user's statistics (admin)
// @Summary Get user statistics
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} UserStatResponse
// @Router /analytics/users/{user_id} [get]
func (h *Handler) GetUserStats(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid user ID")
		return
	}

	stats, err := h.service.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, stats)
}

// GetSystemStats gets system-wide statistics (admin)
// @Summary Get system statistics
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SystemStatResponse
// @Router /admin/analytics/system [get]
func (h *Handler) GetSystemStats(c *gin.Context) {
	stats, err := h.service.GetSystemStats(c.Request.Context())
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, stats)
}

// GetDailyLogins gets daily login counts (admin)
// @Summary Get daily logins
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param start query string true "Start date (YYYY-MM-DD)"
// @Param end query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} TimeSeriesResponse
// @Router /admin/analytics/logins [get]
func (h *Handler) GetDailyLogins(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		utils.BadRequest(c, "start and end dates are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		utils.BadRequest(c, "invalid start date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		utils.BadRequest(c, "invalid end date format")
		return
	}

	data, err := h.service.GetDailyLogins(c.Request.Context(), startDate, endDate.Add(24*time.Hour))
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, data)
}

// GetCourseLeaderboard gets course leaderboard
// @Summary Get course leaderboard
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} LeaderboardResponse
// @Router /analytics/courses/{id}/leaderboard [get]
func (h *Handler) GetCourseLeaderboard(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	leaderboard, err := h.service.GetCourseLeaderboard(c.Request.Context(), courseID, limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, leaderboard)
}
