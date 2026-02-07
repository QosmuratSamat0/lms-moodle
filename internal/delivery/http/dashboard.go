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
