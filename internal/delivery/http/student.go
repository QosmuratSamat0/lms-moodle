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

type CreateStudentRequest struct {
	UserID      string    `json:"user_id" binding:"required"`
	StudentCode string    `json:"student_code" binding:"required"`
	Major       string    `json:"major" binding:"required"`
	Year        int       `json:"year" binding:"required,min=1,max=6"`
	AdmittedAt  time.Time `json:"admitted_at"`
}

type UpdateStudentRequest struct {
	Major  *string  `json:"major"`
	Year   *int     `json:"year"`
	GPA    *float64 `json:"gpa"`
	Status *string  `json:"status"`
}

func NewStudentHandler(service *studentUC.Service) *StudentHandler {
	return &StudentHandler{service: service}
}

// Create creates a new student profile
// @Summary Create student
// @Description Create a new student profile (Admin only)
// @Tags students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateStudentRequest true "Create Student Request"
// @Success 201 {object} student.Student "Created student"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/v1/students [post]
func (h *StudentHandler) Create(c *gin.Context) {
	var req CreateStudentRequest
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

// GetByID returns a student profile by ID
// @Summary Get student by ID
// @Description Returns student profile details by ID
// @Tags students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} student.Student "Student details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Student not found"
// @Router /api/v1/students/{id} [get]
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

// GetByUserID returns a student profile by User ID
// @Summary Get student by User ID
// @Description Returns student profile details by User ID
// @Tags students
// @Security BearerAuth
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} student.Student "Student details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Student not found"
// @Router /api/v1/students/user/{userID} [get]
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

// GetWithDetails returns student profile with full details
// @Summary Get student with details
// @Description Returns student profile along with related user information
// @Tags students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} student.Student "Student details with user info"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/students/{id}/details [get]
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

// GetMyProfile returns current student's profile
// @Summary Get my profile
// @Description Returns the profile of the currently logged-in student
// @Tags students
// @Security BearerAuth
// @Produce json
// @Success 200 {object} student.Student "Student profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/students/me [get]
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

// List returns a list of student profiles
// @Summary List students
// @Description Returns a paginated list of students with optional filtering
// @Tags students
// @Security BearerAuth
// @Produce json
// @Param major query string false "Major"
// @Param year query int false "Year"
// @Param status query string false "Status"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Students list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/students [get]
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

// Update updates student profile details
// @Summary Update student
// @Description Update student profile details like major, year or GPA
// @Tags students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body UpdateStudentRequest true "Update Request"
// @Success 200 {object} student.Student "Updated student"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/students/{id} [put]
func (h *StudentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateStudentRequest
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

// Delete deletes a student profile
// @Summary Delete student
// @Description Deletes a student profile by ID (Admin only)
// @Tags students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/v1/students/{id} [delete]
func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteStudent(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetEnrollments returns student's course enrollments
// @Summary Get student enrollments
// @Description Returns a list of all course enrollments for a specific student
// @Tags students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {array} map[string]interface{} "Enrollments list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/students/{id}/enrollments [get]
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

// GetGroups returns student's group memberships
// @Summary Get student groups
// @Description Returns a list of all groups the student belongs to
// @Tags students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {array} map[string]interface{} "Groups list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/students/{id}/groups [get]
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
