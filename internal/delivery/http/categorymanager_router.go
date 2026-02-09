package http

import (

	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	"github.com/ap1-final-mini-moodle/internal/domain/categorymanager"
	authUC "github.com/ap1-final-mini-moodle/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type CategoryManagerModule struct {
	handler     *CategoryManagerHandler
	authService *authUC.Service
	repo        categorymanager.Repository
}

func NewCategoryManagerModule(handler *CategoryManagerHandler, authService *authUC.Service, repo categorymanager.Repository) *CategoryManagerModule {
	return &CategoryManagerModule{
		handler:     handler,
		authService: authService,
		repo:        repo,
	}
}

func (m *CategoryManagerModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	catMgrs := api.Group("/category-managers")
	catMgrs.Use(middleware.AuthTokenMiddleware(m.authService))
	{
		// VIEW — все авторизованные могут видеть списки и детали
		catMgrs.GET("", m.handler.List)

		catMgrs.GET("/:id", m.handler.GetByID)

		catMgrs.GET("/:id/details", m.handler.GetWithDetails)

		catMgrs.GET("/user-category", m.handler.GetByUserAndCategory)

		catMgrs.GET("/by-user/:userID", m.handler.GetByUserID)

		catMgrs.GET("/by-category/:categoryID", m.handler.GetByCategoryID)

		// UPDATE — требует edit permission
		catMgrs.PUT("/:id",
			middleware.RequireCategoryManagerPermission(m.repo, "edit"),
			m.handler.Update,
		)

		// CREATE — требует admin permission
		catMgrs.POST("",
			middleware.RequireCategoryManagerPermission(m.repo, "admin"),
			m.handler.Create,
		)

		// DELETE — требует admin permission
		catMgrs.DELETE("/:id",
			middleware.RequireCategoryManagerPermission(m.repo, "admin"),
			m.handler.Delete,
		)
	}
}
