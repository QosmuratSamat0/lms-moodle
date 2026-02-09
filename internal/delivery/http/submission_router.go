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
		submissions.POST("", m.handler.Submit)
		submissions.GET("/:id", m.handler.GetByID)
		submissions.GET("/assignment/:assignmentID", m.handler.ListByAssignment)
		submissions.GET("/student/:studentID", m.handler.ListByStudent)
		submissions.DELETE("/:id", m.handler.Delete)
	}
}
