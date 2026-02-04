package http

import (
	"github.com/gin-gonic/gin"
)

type SubmissionModule struct {
	handler *SubmissionHandler
	secret  []byte
}

func NewSubmissionModule(handler *SubmissionHandler, secret []byte) *SubmissionModule {
	return &SubmissionModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *SubmissionModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	submissions := api.Group("/submissions")
	{
		submissions.POST("", m.handler.Submit)
		submissions.GET("/:id", m.handler.GetByID)
		submissions.GET("/assignment/:assignmentID", m.handler.ListByAssignment)
		submissions.GET("/student/:studentID", m.handler.ListByStudent)
		submissions.DELETE("/:id", m.handler.Delete)
	}
}
