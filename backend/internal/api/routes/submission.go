// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/ap1-final-mini-moodle/internal/domain/submission"
	"github.com/gin-gonic/gin"
)

// SubmissionRouter handles submission-related routes
type SubmissionRouter struct {
	submissionHandler *submission.Handler
}

// NewSubmissionRouter creates a new submission router
func NewSubmissionRouter(submissionHandler *submission.Handler) *SubmissionRouter {
	return &SubmissionRouter{
		submissionHandler: submissionHandler,
	}
}

// SetupRoutes configures submission routes
func (sr *SubmissionRouter) SetupRoutes(api *gin.RouterGroup) {
	submissions := api.Group("/submissions")
	{
		submissions.GET("/:id", sr.submissionHandler.GetByID)
		submissions.POST("", sr.submissionHandler.Submit)
		submissions.PUT("/:id", sr.submissionHandler.Update)
		submissions.DELETE("/:id", sr.submissionHandler.Delete)

		// Student submissions
		api.GET("/student/submissions", sr.submissionHandler.ListMySubmissions)

		// Assignment submissions
		api.GET("/assignments/:id/submissions", sr.submissionHandler.ListByAssignment)
	}
}
