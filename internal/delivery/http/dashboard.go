package http

import (
	"net/http"

	dashboardUC "github.com/ap1-final-mini-moodle/internal/usecase/dashboard"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *dashboardUC.Service
}

func NewDashboardHandler(service *dashboardUC.Service) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetStudentDashboard returns dashboard data for a student
// @Summary Get student dashboard
// @Description Gets an overview of courses, assignments, and grades for the student
// @Tags dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dashboard.StudentDashboard
// @Router /api/v1/dashboard/student [get]
func (h *DashboardHandler) GetStudentDashboard(c *gin.Context) {
	userID, _ := c.Get("userID")
	studentID := userID.(string)

	dashboard, err := h.service.GetStudentDashboard(c.Request.Context(), studentID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dashboard)
}

// GetTeacherDashboard returns dashboard data for a teacher
// @Summary Get teacher dashboard
// @Description Gets an overview of courses, students, and grading activities for the teacher
// @Tags dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dashboard.TeacherDashboard
// @Router /api/v1/dashboard/teacher [get]
func (h *DashboardHandler) GetTeacherDashboard(c *gin.Context) {
	userID, _ := c.Get("userID")
	teacherID := userID.(string)

	dashboard, err := h.service.GetTeacherDashboard(c.Request.Context(), teacherID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dashboard)
}

// GetStudentCourseStats returns statistics for a student in a specific course
// @Summary Get student course stats
// @Description Gets detailed statistics like attendance and grades for a specific course
// @Tags dashboard
// @Security BearerAuth
// @Produce json
// @Param courseID path string true "Course ID"
// @Success 200 {object} dashboard.StudentCourseStats
// @Router /api/v1/dashboard/student/courses/{courseID} [get]
func (h *DashboardHandler) GetStudentCourseStats(c *gin.Context) {
	userID, _ := c.Get("userID")
	studentID := userID.(string)
	courseID := c.Param("courseID")

	stats, err := h.service.GetStudentCourseStats(c.Request.Context(), studentID, courseID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
