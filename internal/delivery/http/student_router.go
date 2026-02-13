package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type StudentModule struct {
	handler     *StudentHandler
	authService *authUC.Service
}

func NewStudentModule(handler *StudentHandler, authService *authUC.Service) *StudentModule {
	return &StudentModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *StudentModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	students := api.Group("/students")
	students.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// GET own profile
		students.GET("/me", m.handler.GetMyProfile)

		// VIEW — все могут видеть список и профили студентов
		students.GET("", m.handler.List)
		students.GET("/:id", m.handler.GetByID)
		students.GET("/:id/details", m.handler.GetWithDetails)
		students.GET("/:id/enrollments", m.handler.GetEnrollments)
		students.GET("/:id/groups", m.handler.GetGroups)
		students.GET("/user/:userID", m.handler.GetByUserID)

		// CREATE — только admin (создание профиля студента)
		students.POST("",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Create,
		)

		// UPDATE — student (own) или admin
		students.PUT("/:id",
			m.handler.Update,
		)

		// DELETE — только admin
		students.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)
	}
}
