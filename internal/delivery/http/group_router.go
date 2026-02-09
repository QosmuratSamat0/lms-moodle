package http

import (
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type GroupModule struct {
	handler     *GroupHandler
	authService *authUC.Service
}

func NewGroupModule(handler *GroupHandler, authService *authUC.Service) *GroupModule {
	return &GroupModule{
		handler:     handler,
		authService: authService,
	}
}

func (m *GroupModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthTokenMiddleware(m.authService))

	groups := api.Group("/groups")
	{
		// VIEW — все могут видеть группы
		groups.GET("/:id", m.handler.GetByID)
		groups.GET("/:id/members", m.handler.GetMembers)

		// CREATE — только teacher/admin
		groups.POST("",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Create,
		)

		// UPDATE — teacher/admin
		groups.PUT("/:id",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.Update,
		)

		// DELETE — только admin
		groups.DELETE("/:id",
			middleware.RequireRole("admin", "super_admin"),
			m.handler.Delete,
		)

		// MEMBERS — teacher/admin могут добавлять/удалять
		groups.POST("/:id/members",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.AddMember,
		)
		groups.DELETE("/:id/members/:studentID",
			middleware.RequireRole("teacher", "admin", "super_admin"),
			m.handler.RemoveMember,
		)
	}

	api.GET("/courses/:courseID/groups", m.handler.ListByCourse)
	api.GET("/students/:studentID/groups", m.handler.GetStudentGroups)
}
