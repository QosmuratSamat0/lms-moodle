package http

import (
	"github.com/gin-gonic/gin"
)

type AssignmentModule struct {
	handler *AssignmentHandler
	secret  []byte
}

func NewAssignmentModule(handler *AssignmentHandler, secret []byte) *AssignmentModule {
	return &AssignmentModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *AssignmentModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	assignments := api.Group("/assignments")
	{
		assignments.POST("", m.handler.Create)
		assignments.GET("/:id", m.handler.GetByID)
		assignments.GET("/course/:courseID", m.handler.ListByCourse)
		assignments.DELETE("/:id", m.handler.Delete)
	}
}
