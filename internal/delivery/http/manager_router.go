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
		// GET own profile — manager может видеть свой профиль
		managers.GET("/me", m.handler.GetMyProfile)

		// LIST — только admin может видеть список всех managers
		managers.GET("",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.List,
		)

		// GET by ID — только admin может видеть других managers
		managers.GET("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.GetByID,
		)

		managers.GET("/:id/details",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.GetWithDetails,
		)

		managers.GET("/department",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.GetByDepartment,
		)

		managers.GET("/user/:userID",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.GetByUserID,
		)

		// CREATE — только admin/super_admin
		managers.POST("",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Create,
		)

		// UPDATE — только admin (manager не может обновлять себя через API)
		managers.PUT("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Update,
		)

		// DELETE — только admin/super_admin
		managers.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)

		// MANAGED RESOURCES — только admin может управлять
		managers.POST("/:id/categories",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.AddManagedCategory,
		)
		managers.DELETE("/:id/categories/:categoryID",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.RemoveManagedCategory,
		)
		managers.POST("/:id/teachers",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.AddManagedTeacher,
		)
		managers.DELETE("/:id/teachers/:teacherID",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.RemoveManagedTeacher,
		)
	}
}
