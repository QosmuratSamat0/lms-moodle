package http

import (
	"net/http"
	"strconv"

	"github.com/ap1-final-mini-moodle/internal/domain/teacher"
	teacherUC "github.com/ap1-final-mini-moodle/internal/usecase/teacher"
	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	service *teacherUC.Service
}

func NewTeacherHandler(service *teacherUC.Service) *TeacherHandler {
	return &TeacherHandler{service: service}
}

// Create creates a new teacher profile
// @Summary Create teacher
// @Description Create a new teacher profile (Admin only)
// @Tags teachers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body teacher.CreateTeacherInput true "Create Teacher Request"
// @Success 201 {object} teacher.Teacher "Created teacher"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/v1/teachers [post]
func (h *TeacherHandler) Create(c *gin.Context) {
	var req teacher.CreateTeacherInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.service.CreateTeacher(c.Request.Context(), &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

// GetByID returns a teacher profile by ID
// @Summary Get teacher by ID
// @Description Returns teacher profile details by ID
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID"
// @Success 200 {object} teacher.Teacher "Teacher details"
// @Router /api/v1/teachers/{id} [get]
func (h *TeacherHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	t, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// GetByUserID returns a teacher profile by User ID
// @Summary Get teacher by User ID
// @Description Returns teacher profile details by User ID
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} teacher.Teacher "Teacher details"
// @Router /api/v1/teachers/user/{userID} [get]
func (h *TeacherHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("userID")
	t, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// GetMyProfile returns current teacher's profile
// @Summary Get my profile
// @Description Returns the profile of the currently logged-in teacher
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} teacher.Teacher "Teacher profile"
// @Router /api/v1/teachers/me [get]
func (h *TeacherHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	t, err := h.service.GetByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// List returns a list of teacher profiles
// @Summary List teachers
// @Description Returns a paginated list of teachers
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Teachers list"
// @Router /api/v1/teachers [get]
func (h *TeacherHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	teachers, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   teachers,
		"limit":  limit,
		"offset": offset,
	})
}

// Update updates teacher profile details
// @Summary Update teacher
// @Description Update teacher profile details like degree or department
// @Tags teachers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID"
// @Param request body teacher.UpdateTeacherInput true "Update Request"
// @Success 200 {object} teacher.Teacher "Updated teacher"
// @Router /api/v1/teachers/{id} [put]
func (h *TeacherHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req teacher.UpdateTeacherInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.service.UpdateTeacher(c.Request.Context(), id, &req)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// Delete deletes a teacher profile
// @Summary Delete teacher
// @Description Deletes a teacher profile by ID (Admin only)
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID"
// @Success 204 "No content"
// @Router /api/v1/teachers/{id} [delete]
func (h *TeacherHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTeacher(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetCourses returns teacher's courses
// @Summary Get teacher courses
// @Description Returns a list of all courses taught by a specific teacher
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID"
// @Success 200 {array} map[string]interface{} "Courses list"
// @Router /api/v1/teachers/{id}/courses [get]
func (h *TeacherHandler) GetCourses(c *gin.Context) {
	id := c.Param("id")
	courses, err := h.service.GetTeacherCourses(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// GetGroups returns teacher's groups
// @Summary Get teacher groups
// @Description Returns a list of all groups assigned to a specific teacher
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID"
// @Success 200 {array} map[string]interface{} "Groups list"
// @Router /api/v1/teachers/{id}/groups [get]
func (h *TeacherHandler) GetGroups(c *gin.Context) {
	id := c.Param("id")
	groups, err := h.service.GetTeacherGroups(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}
