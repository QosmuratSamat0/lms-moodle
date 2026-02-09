package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type EnrollmentModule struct {
	handler     *EnrollmentHandler
	authService *authUC.Service
}

func NewEnrollmentModule(handler *EnrollmentHandler, authService *authUC.Service) *EnrollmentModule {
	return &EnrollmentModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *EnrollmentModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	enrollments := api.Group("/enrollments")
	enrollments.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		enrollments.POST("", m.handler.Enroll)
		enrollments.GET("/course/:courseID", m.handler.ListByCourse)
		enrollments.GET("/student/:studentID", m.handler.ListByStudent)
		enrollments.DELETE("/:id", m.handler.Remove)
	}
}
