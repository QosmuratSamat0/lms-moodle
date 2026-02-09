package http

import (
	"github.com/gin-gonic/gin"
)

type CourseCategoryModule struct {
	handler *CourseCategoryHandler
	secret  []byte
}

func NewCourseCategoryModule(handler *CourseCategoryHandler, secret []byte) *CourseCategoryModule {
	return &CourseCategoryModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *CourseCategoryModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	categories := api.Group("/categories")
	{
		// CRUD operations
		categories.POST("", m.handler.Create)
		categories.GET("", m.handler.List)
		categories.GET("/:id", m.handler.GetByID)
		categories.PUT("/:id", m.handler.Update)
		categories.DELETE("/:id", m.handler.Delete)
	}
}
