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

// Create creates a new student group
// @Summary Create group
// @Description Create a new student group in a course (Teacher/Admin only)
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body struct{CourseID string `json:"course_id" binding:"required"`; Name string `json:"name" binding:"required"`; Description string `json:"description"`; MaxStudents int `json:"max_students"`} true "Create Group Request"
// @Success 201 {object} group.Group "Created group"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/v1/groups [post]
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

// GetByID returns a group by ID
// @Summary Get group by ID
// @Description Returns group details by ID
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} group.Group "Group details"
// @Router /api/v1/groups/{id} [get]
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

// ListByCourse returns groups for a specific course
// @Summary List course groups
// @Description Returns a list of all groups in a course
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param courseID path string true "Course ID"
// @Success 200 {array} group.Group "Groups list"
// @Router /api/v1/courses/{courseID}/groups [get]
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

// Update updates group details
// @Summary Update group
// @Description Update group details like name or capacity (Teacher/Admin only)
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body struct{Name *string `json:"name"`; Description *string `json:"description"`; MaxStudents *int `json:"max_students"`} true "Update Request"
// @Success 200 {object} group.Group "Updated group"
// @Router /api/v1/groups/{id} [put]
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

// Delete deletes a group
// @Summary Delete group
// @Description Deletes a group by ID (Admin only)
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Success 204 "No content"
// @Router /api/v1/groups/{id} [delete]
func (h *GroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteGroup(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// AddMember adds a student to a group
// @Summary Add group member
// @Description Adds a student to a specific group (Teacher/Admin only)
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body struct{StudentID string `json:"student_id" binding:"required"`} true "Add Member Request"
// @Success 201 {object} group.GroupMember "Added member"
// @Router /api/v1/groups/{id}/members [post]
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

// RemoveMember removes a student from a group
// @Summary Remove group member
// @Description Removes a student from a specific group (Teacher/Admin only)
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Param studentID path string true "Student ID"
// @Success 204 "No content"
// @Router /api/v1/groups/{id}/members/{studentID} [delete]
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

// GetMembers returns all members of a specific group
// @Summary List group members
// @Description Gets list of students in a group (Course members only)
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {array} group.GroupMember
// @Router /api/v1/groups/{id}/members [get]
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

// GetStudentGroups returns all groups a student is a member of
// @Summary List student's groups
// @Description Gets all groups the student belongs to
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param studentID path string true "Student ID"
// @Success 200 {array} group.Group
// @Router /api/v1/students/{studentID}/groups [get]
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
