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
		// VIEW — все могут видеть
		announcements.GET("/:id", m.handler.GetByID)

		// CREATE — только teacher/admin
		announcements.POST("",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Create,
		)

		// UPDATE — owner (teacher) или admin
		announcements.PUT("/:id",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Update,
		)

		// DELETE — owner или admin
		announcements.DELETE("/:id",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
