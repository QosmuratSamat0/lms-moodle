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
		courses.POST("", m.handler.Create)
		courses.GET("", m.handler.List)
		courses.GET("/:id", m.handler.GetByID)
		courses.PUT("/:id", m.handler.Update)
		courses.DELETE("/:id", m.handler.Delete)
	}
}
