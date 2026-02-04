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

func NewCourseHandler(service *courseUC.Service) *CourseHandler {
	return &CourseHandler{service: service}
}

func (h *CourseHandler) Create(c *gin.Context) {
	var req struct {
		Code        string `json:"code" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		TeacherID   string `json:"teacher_id" binding:"required"`
		MaxPoints   int    `json:"max_points" binding:"required"`
	}
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

func (h *CourseHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	crs, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}
	c.JSON(http.StatusOK, crs)
}

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

func (h *CourseHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		MaxPoints   *int    `json:"max_points"`
		Active      *bool   `json:"active"`
	}
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

func (h *CourseHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
