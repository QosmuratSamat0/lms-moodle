package http

import (
	"github.com/gin-gonic/gin"
)

type CourseModule struct {
	handler *CourseHandler
	secret  []byte
}

func NewCourseModule(handler *CourseHandler, secret []byte) *CourseModule {
	return &CourseModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *CourseModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	courses := api.Group("/courses")
	{
		courses.POST("", m.handler.Create)
		courses.GET("", m.handler.List)
		courses.GET("/:id", m.handler.GetByID)
		courses.PUT("/:id", m.handler.Update)
		courses.DELETE("/:id", m.handler.Delete)
	}
}
