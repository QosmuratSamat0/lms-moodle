package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/group"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
	groupUC "github.com/ap1-final-mini-moodle/internal/usecase/group"
	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	service *groupUC.Service
}

func NewGroupHandler(service *groupUC.Service) *GroupHandler {
	return &GroupHandler{service: service}
}

func (h *GroupHandler) Create(c *gin.Context) {
	var req struct {
		CourseID    string `json:"course_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		MaxStudents int    `json:"max_students"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g, err := h.service.CreateGroup(c.Request.Context(), &group.CreateGroupInput{
		CourseID:    req.CourseID,
		Name:        req.Name,
		Description: req.Description,
		MaxStudents: req.MaxStudents,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *GroupHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	g, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *GroupHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseID")
	groups, err := h.service.GetByCourseID(c.Request.Context(), courseID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

func (h *GroupHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		MaxStudents *int    `json:"max_students"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g, err := h.service.UpdateGroup(c.Request.Context(), id, &group.UpdateGroupInput{
		Name:        req.Name,
		Description: req.Description,
		MaxStudents: req.MaxStudents,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *GroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteGroup(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *GroupHandler) AddMember(c *gin.Context) {
	groupID := c.Param("id")
	var req struct {
		StudentID string `json:"student_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.service.AddMember(c.Request.Context(), &group.AddMemberInput{
		GroupID:   groupID,
		StudentID: req.StudentID,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, member)
}

func (h *GroupHandler) RemoveMember(c *gin.Context) {
	groupID := c.Param("id")
	studentID := c.Param("studentID")

	if err := h.service.RemoveMember(c.Request.Context(), groupID, studentID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *GroupHandler) GetMembers(c *gin.Context) {
	groupID := c.Param("id")
	members, err := h.service.GetMembers(c.Request.Context(), groupID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *GroupHandler) GetStudentGroups(c *gin.Context) {
	studentID := c.Param("studentID")
	groups, err := h.service.GetStudentGroups(c.Request.Context(), studentID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

func getStatusCode(err error) int {
	if appErrors.IsNotFoundError(err) {
		return http.StatusNotFound
	}
	if appErrors.IsForbiddenError(err) {
		return http.StatusForbidden
	}
	if appErrors.IsConflictError(err) {
		return http.StatusConflict
	}
	if appErrors.IsValidationError(err) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
