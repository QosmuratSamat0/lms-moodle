package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AdminModule struct {
	handler     *AdminHandler
	authService *authUC.Service
}

func NewAdminModule(handler *AdminHandler, authService *authUC.Service) *AdminModule {
	return &AdminModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *AdminModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	admins := api.Group("/admins")
	admins.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// Admin & super admin routes
		admins.GET("/me",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.GetMyProfile,
		)

		admins.GET("",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.List,
		)

		admins.GET("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.GetByID,
		)

		admins.GET("/user/:userID",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.GetByUserID,
		)

		admins.PUT("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Update,
		)

		// Only super admin routes
		admins.POST("",
			middleware.RequireRole("super_admin"),
			m.handler.Create,
		)

		admins.DELETE("/:id",
			middleware.RequireRole("super_admin"),
			m.handler.Delete,
		)
	}
}
