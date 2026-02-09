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
		// VIEW — read-only access
		catMgrs.GET("",
			middleware.RequireCategoryManagerPermission(m.repo, "view"),
			m.handler.List,
		)

		catMgrs.GET("/:id",
			middleware.RequireCategoryManagerPermission(m.repo, "view"),
			m.handler.GetByID,
		)

		catMgrs.GET("/:id/details",
			middleware.RequireCategoryManagerPermission(m.repo, "view"),
			m.handler.GetWithDetails,
		)

		catMgrs.GET("/user-category",
			middleware.RequireCategoryManagerPermission(m.repo, "view"),
			m.handler.GetByUserAndCategory,
		)

		catMgrs.GET("/by-user/:userID",
			middleware.RequireCategoryManagerPermission(m.repo, "view"),
			m.handler.GetByUserID,
		)

		catMgrs.GET("/by-category/:categoryID",
			middleware.RequireCategoryManagerPermission(m.repo, "view"),
			m.handler.GetByCategoryID,
		)

		// EDIT — update access
		catMgrs.PUT("/:id",
			middleware.RequireCategoryManagerPermission(m.repo, "edit"),
			m.handler.Update,
		)

		// ADMIN — full access
		catMgrs.POST("",
			middleware.RequireCategoryManagerPermission(m.repo, "admin"),
			m.handler.Create,
		)

		catMgrs.DELETE("/:id",
			middleware.RequireCategoryManagerPermission(m.repo, "admin"),
			m.handler.Delete,
		)
	}
}
