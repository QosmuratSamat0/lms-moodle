package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AnnouncementModule struct {
	handler     *AnnouncementHandler
	authService *authUC.Service
}

func NewAnnouncementModule(handler *AnnouncementHandler, authService *authUC.Service) *AnnouncementModule {
	return &AnnouncementModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *AnnouncementModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthTokenMiddleware(m.authService))

	announcements := api.Group("/announcements")
	{
		announcements.POST("", m.handler.Create)
		announcements.GET("/:id", m.handler.GetByID)
		announcements.PUT("/:id", m.handler.Update)
		announcements.DELETE("/:id", m.handler.Delete)
	}

	api.GET("/courses/:courseID/announcements", m.handler.ListByCourse)
}
