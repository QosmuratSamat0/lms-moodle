package http

import (
	"github.com/gin-gonic/gin"
)

type ManagerModule struct {
	handler *ManagerHandler
	secret  []byte
}

func NewManagerModule(handler *ManagerHandler, secret []byte) *ManagerModule {
	return &ManagerModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *ManagerModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	managers := api.Group("/managers")
	{
		// Get current manager profile (for logged in manager)
		managers.GET("/me", m.handler.GetMyProfile)

		// CRUD operations
		managers.POST("", m.handler.Create)
		managers.GET("", m.handler.List)
		managers.GET("/department", m.handler.GetByDepartment)
		managers.GET("/:id", m.handler.GetByID)
		managers.GET("/:id/details", m.handler.GetWithDetails)
		managers.PUT("/:id", m.handler.Update)
		managers.DELETE("/:id", m.handler.Delete)

		// Managed categories
		managers.POST("/:id/categories", m.handler.AddManagedCategory)
		managers.DELETE("/:id/categories/:categoryID", m.handler.RemoveManagedCategory)

		// Managed teachers
		managers.POST("/:id/teachers", m.handler.AddManagedTeacher)
		managers.DELETE("/:id/teachers/:teacherID", m.handler.RemoveManagedTeacher)

		// Get manager by user ID
		managers.GET("/user/:userID", m.handler.GetByUserID)
	}
}
