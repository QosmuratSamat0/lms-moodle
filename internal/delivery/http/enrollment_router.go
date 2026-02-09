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
		// VIEW — teacher видит в своём курсе, admin видит все
		enrollments.GET("/course/:courseId", m.handler.ListByCourse)
		enrollments.GET("/student/:studentId", m.handler.ListByStudent)

		// CREATE — студент, учитель или админ
		enrollments.POST("",
			middleware.RequireRole("student", "teacher", "admin", "super_admin"),
			m.handler.Enroll,
		)

		// DELETE — только teacher (курса) или admin
		enrollments.DELETE("/:id",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Remove,
		)
	}
}
