package http

import (
	"net/http"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/assignment"
	assignmentUC "github.com/ap1-final-mini-moodle/internal/usecase/assignment"
	"github.com/gin-gonic/gin"
)

type AssignmentHandler struct {
	service *assignmentUC.Service
}

func NewAssignmentHandler(service *assignmentUC.Service) *AssignmentHandler {
	return &AssignmentHandler{service: service}
}

func (h *AssignmentHandler) Create(c *gin.Context) {
	var req struct {
		CourseID    string    `json:"course_id" binding:"required"`
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		MaxPoints   int       `json:"max_points" binding:"required"`
		DueDate     time.Time `json:"due_date" binding:"required"`
	}
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

func (h *AssignmentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "assignment not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

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

func (h *AssignmentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
