package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/course"
	courseUC "github.com/ap1-final-mini-moodle/internal/usecase/course"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service *courseUC.Service
}

type CreateCourseRequest struct {
	Code        string `json:"code" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	TeacherID   string `json:"teacher_id" binding:"required"`
	MaxPoints   int    `json:"max_points" binding:"required"`
}

type UpdateCourseRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	MaxPoints   *int    `json:"max_points"`
	Active      *bool   `json:"active"`
}

func NewCourseHandler(service *courseUC.Service) *CourseHandler {
	return &CourseHandler{service: service}
}

// Create creates a new course
// @Summary Create course
// @Description Create a new course (Admin only)
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateCourseRequest true "Create Course Request"
// @Success 201 {object} course.Course "Created course"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/courses [post]
func (h *CourseHandler) Create(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	crs, err := h.service.Create(&course.CreateCourseInput{
		Code:        req.Code,
		Title:       req.Title,
		Description: req.Description,
		TeacherID:   req.TeacherID,
		MaxPoints:   req.MaxPoints,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, crs)
}

// GetByID returns a course by ID
// @Summary Get course by ID
// @Description Returns course details by ID
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} course.Course "Course details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Course not found"
// @Router /api/v1/courses/{id} [get]
func (h *CourseHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	crs, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}
	c.JSON(http.StatusOK, crs)
}

// List returns a list of courses
// @Summary List courses
// @Description Returns a paginated list of courses
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} course.Course "Courses list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/courses [get]
func (h *CourseHandler) List(c *gin.Context) {
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	courses, err := h.service.List(req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// Update updates course details
// @Summary Update course
// @Description Update course details (Teacher/Admin only)
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param request body UpdateCourseRequest true "Update Request"
// @Success 200 {object} course.Course "Updated course"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Course not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/courses/{id} [put]
func (h *CourseHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	crs, err := h.service.Update(id, &course.UpdateCourseInput{
		Title:       req.Title,
		Description: req.Description,
		MaxPoints:   req.MaxPoints,
		Active:      req.Active,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, crs)
}

// Delete deletes a course
// @Summary Delete course
// @Description Deletes a course by ID (Admin only)
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Course not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/courses/{id} [delete]
func (h *CourseHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
