package http

import (
	"github.com/gin-gonic/gin"
)

type GradeModule struct {
	handler *GradeHandler
	secret  []byte
}

func NewGradeModule(handler *GradeHandler, secret []byte) *GradeModule {
	return &GradeModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *GradeModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	grades := api.Group("/grades")
	{
		grades.POST("", m.handler.Grade)
		grades.GET("/:id", m.handler.GetByID)
		grades.DELETE("/:id", m.handler.Delete)
	}
}
