package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type TeacherModule struct {
	handler     *TeacherHandler
	authService *authUC.Service
}

func NewTeacherModule(handler *TeacherHandler, authService *authUC.Service) *TeacherModule {
	return &TeacherModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *TeacherModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	teachers := api.Group("/teachers")
	teachers.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// GET own profile
		teachers.GET("/me", m.handler.GetMyProfile)

		// VIEW — все могут видеть список и профили учителей
		teachers.GET("", m.handler.List)
		teachers.GET("/:id", m.handler.GetByID)
		teachers.GET("/:id/courses", m.handler.GetCourses)
		teachers.GET("/:id/groups", m.handler.GetGroups)
		teachers.GET("/user/:userID", m.handler.GetByUserID)

		// CREATE — только admin
		teachers.POST("",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Create,
		)

		// UPDATE — teacher (own) или admin
		teachers.PUT("/:id",
			m.handler.Update,
		)

		// DELETE — только admin
		teachers.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
