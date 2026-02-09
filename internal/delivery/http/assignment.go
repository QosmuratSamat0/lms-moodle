package http

import (
	"net/http"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/assignment"
	assignmentUC "github.com/ap1-final-mini-moodle/internal/usecase/assignment"
	"github.com/gin-gonic/gin"
)

type CreateAssignmentRequest struct {
	CourseID    string    `json:"course_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	MaxPoints   int       `json:"max_points" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
}

type AssignmentHandler struct {
	service *assignmentUC.Service
}

func NewAssignmentHandler(service *assignmentUC.Service) *AssignmentHandler {
	return &AssignmentHandler{service: service}
}

// Create creates a new assignment
// @Summary Create assignment
// @Description Create a new course assignment (Teacher/Admin only)
// @Tags assignments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateAssignmentRequest true "Create Assignment Request"
// @Success 201 {object} assignment.Assignment "Created assignment"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/assignments [post]
func (h *AssignmentHandler) Create(c *gin.Context) {
	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.service.Create(&assignment.CreateAssignmentInput{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		MaxPoints:   req.MaxPoints,
		DueDate:     req.DueDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

// GetByID returns an assignment by ID
// @Summary Get assignment by ID
// @Description Returns assignment details by ID
// @Tags assignments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} assignment.Assignment "Assignment details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Assignment not found"
// @Router /api/v1/assignments/{id} [get]
func (h *AssignmentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "assignment not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// ListByCourse returns assignments for a specific course
// @Summary List course assignments
// @Description Returns a paginated list of assignments for a course
// @Tags assignments
// @Security BearerAuth
// @Produce json
// @Param courseID path string true "Course ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} assignment.Assignment "Assignments list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/assignments/course/{courseID} [get]
func (h *AssignmentHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assignments, err := h.service.ListByCourse(courseID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

// Delete deletes an assignment
// @Summary Delete assignment
// @Description Deletes an assignment by ID (Teacher/Admin only)
// @Tags assignments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Assignment not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/assignments/{id} [delete]
func (h *AssignmentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
