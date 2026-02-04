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

func NewAttendanceHandler(service *attendanceUC.Service) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) Record(c *gin.Context) {
	var req struct {
		CourseID  string    `json:"course_id" binding:"required"`
		StudentID string    `json:"student_id" binding:"required"`
		Date      time.Time `json:"date" binding:"required"`
		Present   bool      `json:"present"`
	}
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

func (h *AttendanceHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseID")
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

func (h *AttendanceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
