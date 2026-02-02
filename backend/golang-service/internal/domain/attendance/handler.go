package attendance

import (
	"net/http"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles attendance-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new attendance handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreateSession creates a new attendance session
// @Summary Create attendance session
// @Tags attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateSessionRequest true "Session details"
// @Success 201 {object} SessionResponse
// @Failure 400 {object} utils.Response
// @Router /attendance/sessions [post]
func (h *Handler) CreateSession(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	session, err := h.service.CreateSession(c.Request.Context(), userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, session)
}

// GetSession retrieves a session with attendance marks
// @Summary Get attendance session
// @Tags attendance
// @Security BearerAuth
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} SessionAttendanceResponse
// @Failure 404 {object} utils.Response
// @Router /attendance/sessions/{id} [get]
func (h *Handler) GetSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid session ID")
		return
	}

	session, err := h.service.GetSession(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, session)
}

// UpdateSession updates an attendance session
// @Summary Update attendance session
// @Tags attendance
// @Security BearerAuth
// @Accept json
// @Param id path string true "Session ID"
// @Param request body UpdateSessionRequest true "Update details"
// @Success 200 {object} utils.Response
// @Failure 400,403,404 {object} utils.Response
// @Router /attendance/sessions/{id} [put]
func (h *Handler) UpdateSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid session ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	var req UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.UpdateSession(c.Request.Context(), id, userID, role, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "session updated successfully"})
}

// DeleteSession deletes an attendance session
// @Summary Delete attendance session
// @Tags attendance
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 204
// @Failure 403,404 {object} utils.Response
// @Router /attendance/sessions/{id} [delete]
func (h *Handler) DeleteSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid session ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	if err := h.service.DeleteSession(c.Request.Context(), id, userID, role); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSessionsByCourse lists sessions for a course
// @Summary List course attendance sessions
// @Tags attendance
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} SessionListResponse
// @Router /courses/{id}/attendance/sessions [get]
func (h *Handler) ListSessionsByCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListSessionsByCourse(c.Request.Context(), courseID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// MarkAttendance marks attendance for a student
// @Summary Mark student attendance
// @Tags attendance
// @Security BearerAuth
// @Accept json
// @Param session_id path string true "Session ID"
// @Param request body MarkAttendanceRequest true "Attendance mark"
// @Success 200 {object} utils.Response
// @Failure 400,404 {object} utils.Response
// @Router /attendance/sessions/{session_id}/marks [post]
func (h *Handler) MarkAttendance(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid session ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.MarkAttendance(c.Request.Context(), sessionID, userID, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "attendance marked successfully"})
}

// BulkMarkAttendance marks attendance for multiple students
// @Summary Bulk mark attendance
// @Tags attendance
// @Security BearerAuth
// @Accept json
// @Param session_id path string true "Session ID"
// @Param request body BulkMarkAttendanceRequest true "Bulk marks"
// @Success 200 {object} utils.Response
// @Failure 400,404 {object} utils.Response
// @Router /attendance/sessions/{session_id}/marks/bulk [post]
func (h *Handler) BulkMarkAttendance(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid session ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req BulkMarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.BulkMarkAttendance(c.Request.Context(), sessionID, userID, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "attendance marked for all students successfully"})
}

// GetStudentCourseSummary gets attendance summary for a student in a course
// @Summary Get student attendance summary
// @Tags attendance
// @Security BearerAuth
// @Produce json
// @Param student_id path string true "Student ID"
// @Param course_id path string true "Course ID"
// @Success 200 {object} AttendanceSummaryResponse
// @Router /attendance/students/{student_id}/courses/{course_id}/summary [get]
func (h *Handler) GetStudentCourseSummary(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid student ID")
		return
	}

	courseIDStr := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	summary, err := h.service.GetStudentCourseSummary(c.Request.Context(), studentID, courseID)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, summary)
}
