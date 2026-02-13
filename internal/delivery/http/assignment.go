package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/assignment"
	assignmentUC "github.com/ap1-final-mini-moodle/internal/usecase/assignment"
	enrollmentUC "github.com/ap1-final-mini-moodle/internal/usecase/enrollment"
	notificationUC "github.com/ap1-final-mini-moodle/internal/usecase/notification"
	"github.com/gin-gonic/gin"
)

type CreateAssignmentRequest struct {
	CourseID         string   `json:"course_id" binding:"required"`
	Title            string   `json:"title" binding:"required"`
	Description      string   `json:"description"`
	MaxPoints        float64  `json:"max_points"`
	DueAt            *string  `json:"due_at"`
	AllowLate        bool     `json:"allow_late"`
	GradingCategory  string   `json:"grading_category"`
	WeightPercentage float64  `json:"weight_percentage"`
	FileURL          *string  `json:"file_url"`
}

type UpdateAssignmentRequest struct {
	Title            *string  `json:"title"`
	Description      *string  `json:"description"`
	MaxPoints        *float64 `json:"max_points"`
	DueAt            *string  `json:"due_at"`
	AllowLate        *bool    `json:"allow_late"`
	GradingCategory  *string  `json:"grading_category"`
	WeightPercentage *float64 `json:"weight_percentage"`
	FileURL          *string  `json:"file_url"`
}

type AssignmentHandler struct {
	service     *assignmentUC.Service
	notifSvc    *notificationUC.Service
	enrollSvc   *enrollmentUC.Service
}

func NewAssignmentHandler(service *assignmentUC.Service, notifSvc *notificationUC.Service, enrollSvc *enrollmentUC.Service) *AssignmentHandler {
	return &AssignmentHandler{service: service, notifSvc: notifSvc, enrollSvc: enrollSvc}
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

	var dueAt *time.Time
	if req.DueAt != nil && *req.DueAt != "" {
		raw := *req.DueAt
		// Try multiple date formats: RFC3339, datetime-local, date-only
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02T15:04",
			"2006-01-02",
		}
		for _, f := range formats {
			parsed, err := time.Parse(f, raw)
			if err == nil {
				dueAt = &parsed
				break
			}
		}
	}

	// Get teacher ID from JWT claims if available
	var teacherID *string = nil
	if uid, exists := c.Get("userID"); exists {
		uidStr := uid.(string)
		teacherID = &uidStr
	}

	gradingCategory := req.GradingCategory
	if gradingCategory == "" {
		gradingCategory = "register_midterm"
	}

	a, err := h.service.Create(&assignment.CreateAssignmentInput{
		CourseID:           req.CourseID,
		Title:              req.Title,
		Description:        req.Description,
		MaxPoints:          req.MaxPoints,
		DueAt:              dueAt,
		AllowLate:          req.AllowLate,
		CreatedByTeacherID: teacherID,
		GradingCategory:    gradingCategory,
		WeightPercentage:   req.WeightPercentage,
		FileURL:            req.FileURL,
	})
	if err != nil {
		log.Printf("[ASSIGNMENT] Create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notify all enrolled students about the new assignment
	if h.notifSvc != nil && h.enrollSvc != nil {
		go func() {
			enrollments, err := h.enrollSvc.ListByCourse(a.CourseID, 0, 1000)
			if err != nil {
				log.Printf("[NOTIFICATION] Failed to get enrollments: %v", err)
				return
			}
			var studentIDs []string
			for _, e := range enrollments {
				if e.Status == "active" {
					studentIDs = append(studentIDs, e.StudentID)
				}
			}
			if len(studentIDs) > 0 {
				title := fmt.Sprintf("New Assignment: %s", a.Title)
				msg := fmt.Sprintf("A new assignment \"%s\" has been posted in your course.", a.Title)
				if a.DueAt != nil {
					msg += fmt.Sprintf(" Due: %s", a.DueAt.Format("Jan 2, 2006 15:04"))
				}
				h.notifSvc.NotifyMany(studentIDs, "info", title, msg)
			}
		}()
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
// @Param courseId path string true "Course ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} assignment.Assignment "Assignments list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/assignments/course/{courseId} [get]
func (h *AssignmentHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseId")
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
	c.Status(http.StatusNoContent)
}

// Update updates an existing assignment
// @Summary Update assignment
// @Description Update an existing assignment by ID (Teacher/Admin only)
// @Tags assignments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Param request body UpdateAssignmentRequest true "Update Assignment Request"
// @Success 200 {object} assignment.Assignment "Updated assignment"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Assignment not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/assignments/{id} [put]
func (h *AssignmentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := &assignment.UpdateAssignmentInput{
		Title:            req.Title,
		Description:      req.Description,
		MaxPoints:        req.MaxPoints,
		AllowLate:        req.AllowLate,
		GradingCategory:  req.GradingCategory,
		WeightPercentage: req.WeightPercentage,
		FileURL:          req.FileURL,
	}

	// Parse due_at with multiple date formats
	if req.DueAt != nil && *req.DueAt != "" {
		raw := *req.DueAt
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02T15:04",
			"2006-01-02",
		}
		for _, f := range formats {
			parsed, err := time.Parse(f, raw)
			if err == nil {
				input.DueAt = &parsed
				break
			}
		}
	}

	a, err := h.service.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}
