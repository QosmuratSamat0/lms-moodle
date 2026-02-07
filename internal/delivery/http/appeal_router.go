package http

import (
	"github.com/gin-gonic/gin"
)

type AppealModule struct {
	handler *AppealHandler
	secret  []byte
}

func NewAppealModule(handler *AppealHandler, secret []byte) *AppealModule {
	return &AppealModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *AppealModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	appeals := api.Group("/appeals")
	{
		// Student routes
		appeals.POST("", m.handler.Create)
		appeals.GET("", m.handler.GetStudentAppeals)
		appeals.GET("/:id", m.handler.GetByID)
		appeals.DELETE("/:id", m.handler.Delete)

		// Teacher route - resolve appeal
		appeals.PUT("/:id/resolve", m.handler.Resolve)
	}

	// Teacher-specific appeal routes
	api.GET("/teacher/appeals", m.handler.GetTeacherAppeals)
}
