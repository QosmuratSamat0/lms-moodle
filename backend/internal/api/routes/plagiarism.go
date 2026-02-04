// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/plagiarism"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

// PlagiarismRouter handles plagiarism-related routes
type PlagiarismRouter struct {
	plagiarismHandler *plagiarism.Handler
}

// NewPlagiarismRouter creates a new plagiarism router
func NewPlagiarismRouter(plagiarismHandler *plagiarism.Handler) *PlagiarismRouter {
	return &PlagiarismRouter{
		plagiarismHandler: plagiarismHandler,
	}
}

// SetupRoutes configures plagiarism routes
func (pr *PlagiarismRouter) SetupRoutes(api *gin.RouterGroup) {
	plagiarismRoutes := api.Group("/plagiarism")
	plagiarismRoutes.Use(middleware.RequireRole("teacher", "admin"))
	{
		plagiarismRoutes.POST("/check", pr.plagiarismHandler.CheckSubmission)
		plagiarismRoutes.GET("/reports/:id", pr.plagiarismHandler.GetReport)
		plagiarismRoutes.GET("/submissions/:submission_id", pr.plagiarismHandler.GetReportBySubmission)
		plagiarismRoutes.GET("/assignments/:assignment_id", pr.plagiarismHandler.ListByAssignment)
		plagiarismRoutes.GET("/suspicious", pr.plagiarismHandler.ListSuspicious)
		plagiarismRoutes.POST("/submissions/:submission_id/recheck", pr.plagiarismHandler.Recheck)
	}
}
