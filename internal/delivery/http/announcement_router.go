package http

import (
	"github.com/gin-gonic/gin"
)

type AnnouncementModule struct {
	handler *AnnouncementHandler
	secret  []byte
}

func NewAnnouncementModule(handler *AnnouncementHandler, secret []byte) *AnnouncementModule {
	return &AnnouncementModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *AnnouncementModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	announcements := api.Group("/announcements")
	{
		announcements.POST("", m.handler.Create)
		announcements.GET("/:id", m.handler.GetByID)
		announcements.PUT("/:id", m.handler.Update)
		announcements.DELETE("/:id", m.handler.Delete)
	}

	// Course-specific announcement routes
	api.GET("/courses/:courseID/announcements", m.handler.ListByCourse)
}
