package http

import (
	"github.com/gin-gonic/gin"
)

type DashboardModule struct {
	handler *DashboardHandler
	secret  []byte
}

func NewDashboardModule(handler *DashboardHandler, secret []byte) *DashboardModule {
	return &DashboardModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *DashboardModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	dashboard := api.Group("/dashboard")
	{
		// Student dashboard
		dashboard.GET("/student", m.handler.GetStudentDashboard)
		dashboard.GET("/student/courses/:courseID", m.handler.GetStudentCourseStats)

		// Teacher dashboard
		dashboard.GET("/teacher", m.handler.GetTeacherDashboard)
	}
}
