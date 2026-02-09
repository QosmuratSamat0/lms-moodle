package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type GroupModule struct {
	handler     *GroupHandler
	authService *authUC.Service
}

func NewGroupModule(handler *GroupHandler, authService *authUC.Service) *GroupModule {
	return &GroupModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *GroupModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthTokenMiddleware(m.authService))

	groups := api.Group("/groups")
	{
		groups.POST("", m.handler.Create)
		groups.GET("/:id", m.handler.GetByID)
		groups.PUT("/:id", m.handler.Update)
		groups.DELETE("/:id", m.handler.Delete)
		groups.POST("/:id/members", m.handler.AddMember)
		groups.DELETE("/:id/members/:studentID", m.handler.RemoveMember)
		groups.GET("/:id/members", m.handler.GetMembers)
	}

	api.GET("/courses/:courseID/groups", m.handler.ListByCourse)
	api.GET("/students/:studentID/groups", m.handler.GetStudentGroups)
}
