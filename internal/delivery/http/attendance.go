package http

import (
	"log"
	"net/http"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/attendance"
	attendanceUC "github.com/ap1-final-mini-moodle/internal/usecase/attendance"
	notificationUC "github.com/ap1-final-mini-moodle/internal/usecase/notification"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service     *attendanceUC.Service
	notifSvc    *notificationUC.Service
}

func NewAttendanceHandler(service *attendanceUC.Service, notifSvc *notificationUC.Service) *AttendanceHandler {
	return &AttendanceHandler{service: service, notifSvc: notifSvc}
}

// --- Session endpoints ---

type CreateSessionRequest struct {
	CourseID  string `json:"course_id" binding:"required"`
	Date      string `json:"date" binding:"required"`      // "2026-02-13"
	StartTime string `json:"start_time"`                   // "16:00"
	EndTime   string `json:"end_time"`                     // "16:50"
}

func (h *AttendanceHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date, err = time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
	}

	// Combine date + start_time
	startsAt := date
	if req.StartTime != "" {
		t, err := time.Parse("15:04", req.StartTime)
		if err == nil {
			startsAt = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
		}
	}

	// Combine date + end_time
	var endsAt *time.Time
	if req.EndTime != "" {
		t, err := time.Parse("15:04", req.EndTime)
		if err == nil {
			end := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
			endsAt = &end
		}
	}

	teacherID := ""
	if uid, exists := c.Get("userID"); exists {
		teacherID = uid.(string)
	}

	session, err := h.service.CreateSession(&attendance.CreateSessionInput{
		CourseID:  req.CourseID,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		TeacherID: teacherID,
	})
	if err != nil {
		log.Printf("[ATTENDANCE] Create session error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, session)
}

func (h *AttendanceHandler) GetSession(c *gin.Context) {
	id := c.Param("id")
	session, err := h.service.GetSession(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *AttendanceHandler) ListSessionsByCourse(c *gin.Context) {
	courseID := c.Param("courseId")
	sessions, err := h.service.ListSessionsByCourse(courseID)
	if err != nil {
		log.Printf("[ATTENDANCE] List sessions error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sessions == nil {
		sessions = []*attendance.AttendanceSession{}
	}
	c.JSON(http.StatusOK, gin.H{"sessions": sessions, "total": len(sessions)})
}

func (h *AttendanceHandler) DeleteSession(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteSession(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Mark endpoints ---

type BulkMarkRequest struct {
	SessionID string                  `json:"session_id" binding:"required"`
	Marks     []attendance.SingleMark `json:"marks" binding:"required"`
}

func (h *AttendanceHandler) BulkMark(c *gin.Context) {
	var req BulkMarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.BulkMark(&attendance.BulkMarkInput{
		SessionID: req.SessionID,
		Marks:     req.Marks,
	})
	if err != nil {
		log.Printf("[ATTENDANCE] Bulk mark error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notify each student about their attendance
	if h.notifSvc != nil {
		go func() {
			for _, m := range req.Marks {
				title := "Attendance Recorded"
				msg := "Your attendance has been marked as: " + m.Status
				h.notifSvc.NotifyMany([]string{m.StudentID}, "info", title, msg)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance marked successfully"})
}

func (h *AttendanceHandler) GetMarksBySession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	marks, err := h.service.GetMarksBySession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if marks == nil {
		marks = []*attendance.AttendanceMark{}
	}
	c.JSON(http.StatusOK, gin.H{"marks": marks, "total": len(marks)})
}

// Student endpoints
func (h *AttendanceHandler) GetMyAttendance(c *gin.Context) {
	courseID := c.Param("courseId")
	studentID := ""
	if uid, exists := c.Get("userID"); exists {
		studentID = uid.(string)
	}

	marks, err := h.service.GetStudentAttendance(studentID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if marks == nil {
		marks = []*attendance.AttendanceMark{}
	}

	summary, _ := h.service.GetStudentSummary(studentID, courseID)

	c.JSON(http.StatusOK, gin.H{
		"marks":   marks,
		"summary": summary,
	})
}

func (h *AttendanceHandler) GetStudentSummary(c *gin.Context) {
	courseID := c.Param("courseId")
	studentID := c.Param("studentId")

	summary, err := h.service.GetStudentSummary(studentID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
