package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type ManagerModule struct {
	handler     *ManagerHandler
	authService *authUC.Service
}

func NewManagerModule(handler *ManagerHandler, authService *authUC.Service) *ManagerModule {
	return &ManagerModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *ManagerModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	managers := api.Group("/managers")
	managers.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		managers.GET("/me", m.handler.GetMyProfile)
		managers.POST("", m.handler.Create)
		managers.GET("", m.handler.List)
		managers.GET("/department", m.handler.GetByDepartment)
		managers.GET("/:id", m.handler.GetByID)
		managers.GET("/:id/details", m.handler.GetWithDetails)
		managers.PUT("/:id", m.handler.Update)
		managers.DELETE("/:id", m.handler.Delete)
		managers.POST("/:id/categories", m.handler.AddManagedCategory)
		managers.DELETE("/:id/categories/:categoryID", m.handler.RemoveManagedCategory)
		managers.POST("/:id/teachers", m.handler.AddManagedTeacher)
		managers.DELETE("/:id/teachers/:teacherID", m.handler.RemoveManagedTeacher)
		managers.GET("/user/:userID", m.handler.GetByUserID)
	}
}
