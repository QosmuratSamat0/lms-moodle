package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AssignmentModule struct {
	handler     *AssignmentHandler
	authService *authUC.Service
}

func NewAssignmentModule(handler *AssignmentHandler, authService *authUC.Service) *AssignmentModule {
	return &AssignmentModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *AssignmentModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	assignments := api.Group("/assignments")
	assignments.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// VIEW — студенты видят свои, учителя свои, админы все
		assignments.GET("/:id", m.handler.GetByID)
		assignments.GET("/course/:courseID", m.handler.ListByCourse)

		// CREATE — только teacher/admin
		assignments.POST("",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Create,
		)

		// DELETE — только teacher (owner) или admin
		assignments.DELETE("/:id",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
