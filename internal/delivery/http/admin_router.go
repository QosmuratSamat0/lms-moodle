package http

import (
	"github.com/gin-gonic/gin"
)

type AdminModule struct {
	handler *AdminHandler
	secret  []byte
}

func NewAdminModule(handler *AdminHandler, secret []byte) *AdminModule {
	return &AdminModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *AdminModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	admins := api.Group("/admins")
	{
		// Get current admin profile (for logged in admin)
		admins.GET("/me", m.handler.GetMyProfile)

		// CRUD operations
		admins.POST("", m.handler.Create)
		admins.GET("", m.handler.List)
		admins.GET("/:id", m.handler.GetByID)
		admins.PUT("/:id", m.handler.Update)
		admins.DELETE("/:id", m.handler.Delete)

		// Get admin by user ID
		admins.GET("/user/:userID", m.handler.GetByUserID)
	}
}
