package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/enrollment"
	enrollmentUC "github.com/ap1-final-mini-moodle/internal/usecase/enrollment"
	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	service *enrollmentUC.Service
}

func NewEnrollmentHandler(service *enrollmentUC.Service) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}

func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	var req struct {
		CourseID  string `json:"course_id" binding:"required"`
		StudentID string `json:"student_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e, err := h.service.Enroll(&enrollment.CreateEnrollmentInput{
		CourseID:  req.CourseID,
		StudentID: req.StudentID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *EnrollmentHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	enrollments, err := h.service.ListByCourse(courseID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}

func (h *EnrollmentHandler) ListByStudent(c *gin.Context) {
	studentID := c.Param("studentID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	enrollments, err := h.service.ListByStudent(studentID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}

func (h *EnrollmentHandler) Remove(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Remove(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
