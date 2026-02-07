package http

import (
	"github.com/gin-gonic/gin"
)

type GroupModule struct {
	handler *GroupHandler
	secret  []byte
}

func NewGroupModule(handler *GroupHandler, secret []byte) *GroupModule {
	return &GroupModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *GroupModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

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

	// Course-specific group routes
	api.GET("/courses/:courseID/groups", m.handler.ListByCourse)

	// Student-specific group routes
	api.GET("/students/:studentID/groups", m.handler.GetStudentGroups)
}
