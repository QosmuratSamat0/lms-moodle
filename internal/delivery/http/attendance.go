package http

import (
	"net/http"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/attendance"
	attendanceUC "github.com/ap1-final-mini-moodle/internal/usecase/attendance"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service *attendanceUC.Service
}

type CreateAttendanceRequest struct {
	CourseID  string    `json:"course_id" binding:"required"`
	StudentID string    `json:"student_id" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	Present   bool      `json:"present"`
}

func NewAttendanceHandler(service *attendanceUC.Service) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

// Record records attendance for a student
// @Summary Record attendance
// @Description Records attendance for a student in a course (Teacher/Admin only)
// @Tags attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateAttendanceRequest true "Attendance Request"
// @Success 201 {object} attendance.Attendance "Created attendance"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/attendance [post]
func (h *AttendanceHandler) Record(c *gin.Context) {
	var req CreateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.service.Record(&attendance.CreateAttendanceInput{
		CourseID:  req.CourseID,
		StudentID: req.StudentID,
		Date:      req.Date,
		Present:   req.Present,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

// ListByCourse returns attendance records for a specific course
// @Summary List course attendance
// @Description Returns a paginated list of attendance records for a course
// @Tags attendance
// @Security BearerAuth
// @Produce json
// @Param courseId path string true "Course ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} attendance.Attendance "Attendance list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/attendance/course/{courseId} [get]
func (h *AttendanceHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseId")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	attendances, err := h.service.ListByCourse(courseID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attendances)
}

// Delete deletes an attendance record
// @Summary Delete attendance
// @Description Deletes an attendance record by ID (Admin only)
// @Tags attendance
// @Security BearerAuth
// @Produce json
// @Param id path string true "Attendance ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Attendance not found"
// @Router /api/v1/attendance/{id} [delete]
func (h *AttendanceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
