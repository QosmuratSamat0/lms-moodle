package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type CourseModule struct {
	handler     *CourseHandler
	authService *authUC.Service
}

func NewCourseModule(handler *CourseHandler, authService *authUC.Service) *CourseModule {
	return &CourseModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *CourseModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	courses := api.Group("/courses")
	courses.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// VIEW — все авторизованные
		courses.GET("", m.handler.List)
		courses.GET("/:id", m.handler.GetByID)

		// CREATE — только admin/super_admin
		courses.POST("",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Create,
		)

		// UPDATE — teacher (owner) или admin
		courses.PUT("/:id",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Update,
		)

		// DELETE — только admin/super_admin
		courses.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
