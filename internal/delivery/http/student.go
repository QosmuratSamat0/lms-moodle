package http

import (
	"net/http"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/student"
	studentUC "github.com/ap1-final-mini-moodle/internal/usecase/student"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	service *studentUC.Service
}

func NewStudentHandler(service *studentUC.Service) *StudentHandler {
	return &StudentHandler{service: service}
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req struct {
		UserID      string    `json:"user_id" binding:"required"`
		StudentCode string    `json:"student_code" binding:"required"`
		Major       string    `json:"major" binding:"required"`
		Year        int       `json:"year" binding:"required,min=1,max=6"`
		AdmittedAt  time.Time `json:"admitted_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := h.service.CreateStudent(c.Request.Context(), &student.CreateStudentInput{
		UserID:      req.UserID,
		StudentCode: req.StudentCode,
		Major:       req.Major,
		Year:        req.Year,
		AdmittedAt:  req.AdmittedAt,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	s, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *StudentHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("userID")
	s, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *StudentHandler) GetWithDetails(c *gin.Context) {
	id := c.Param("id")
	details, err := h.service.GetWithDetails(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, details)
}

func (h *StudentHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	s, err := h.service.GetByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *StudentHandler) List(c *gin.Context) {
	var filter student.StudentFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	students, total, err := h.service.List(c.Request.Context(), &filter)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   students,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *StudentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Major  *string  `json:"major"`
		Year   *int     `json:"year"`
		GPA    *float64 `json:"gpa"`
		Status *string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := h.service.UpdateStudent(c.Request.Context(), id, &student.UpdateStudentInput{
		Major:  req.Major,
		Year:   req.Year,
		GPA:    req.GPA,
		Status: req.Status,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteStudent(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *StudentHandler) GetEnrollments(c *gin.Context) {
	id := c.Param("id")
	enrollments, err := h.service.GetEnrollments(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}

func (h *StudentHandler) GetGroups(c *gin.Context) {
	id := c.Param("id")
	groups, err := h.service.GetGroups(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}
