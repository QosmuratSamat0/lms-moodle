package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type SubmissionModule struct {
	handler     *SubmissionHandler
	authService *authUC.Service
}

func NewSubmissionModule(handler *SubmissionHandler, authService *authUC.Service) *SubmissionModule {
	return &SubmissionModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *SubmissionModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	submissions := api.Group("/submissions")
	submissions.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// VIEW — студент видит свои, учитель видит все в своём курсе, админ видит все
		submissions.GET("/:id", m.handler.GetByID)
		submissions.GET("/assignment/:assignmentId", m.handler.ListByAssignment)
		submissions.GET("/student/:studentId", m.handler.ListByStudent)

		// CREATE — только student/teacher/admin
		submissions.POST("",
			middleware.RequireRole("student", "teacher", "admin", "super_admin"),
			m.handler.Submit,
		)

		// DELETE — только admin или owner
		submissions.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin", "teacher"),
			m.handler.Delete,
		)
	}
}
