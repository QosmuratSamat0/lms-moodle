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

// Enroll enrolls a student in a course
// @Summary Enroll student
// @Description Enrolls a student in a course
// @Tags enrollments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body struct{CourseID string `json:"course_id" binding:"required"`; StudentID string `json:"student_id" binding:"required"`} true "Enrollment Request"
// @Success 201 {object} enrollment.Enrollment "Created enrollment"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/enrollments [post]
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

// ListByCourse returns enrollments for a specific course
// @Summary List course enrollments
// @Description Returns a paginated list of student enrollments for a course
// @Tags enrollments
// @Security BearerAuth
// @Produce json
// @Param courseID path string true "Course ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} enrollment.Enrollment "Enrollments list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/enrollments/course/{courseID} [get]
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

// ListByStudent returns enrollments for a specific student
// @Summary List student enrollments
// @Description Returns a paginated list of course enrollments for a student
// @Tags enrollments
// @Security BearerAuth
// @Produce json
// @Param studentID path string true "Student ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} enrollment.Enrollment "Enrollments list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/enrollments/student/{studentID} [get]
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

// Remove removes a student enrollment
// @Summary Remove enrollment
// @Description Deletes an enrollment by ID (Teacher/Admin only)
// @Tags enrollments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/enrollments/{id} [delete]
func (h *EnrollmentHandler) Remove(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Remove(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
