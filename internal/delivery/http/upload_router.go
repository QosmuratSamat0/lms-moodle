package http

import (
	"github.com/gin-gonic/gin"
)

type UploadModule struct {
	handler *UploadHandler
	secret  []byte
}

func NewUploadModule(handler *UploadHandler, secret []byte) *UploadModule {
	return &UploadModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *UploadModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	uploads := api.Group("/uploads")
	{
		uploads.POST("", m.handler.Upload)
		uploads.GET("/:id", m.handler.GetByID)
		uploads.GET("/user/:userID", m.handler.ListByUser)
		uploads.DELETE("/:id", m.handler.Delete)
	}
}
