package http

import (
	"github.com/gin-gonic/gin"
)

type CategoryManagerModule struct {
	handler *CategoryManagerHandler
	secret  []byte
}

func NewCategoryManagerModule(handler *CategoryManagerHandler, secret []byte) *CategoryManagerModule {
	return &CategoryManagerModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *CategoryManagerModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	catMgrs := api.Group("/category-managers")
	{
		// CRUD operations
		catMgrs.POST("", m.handler.Create)
		catMgrs.GET("", m.handler.List)
		catMgrs.GET("/:id", m.handler.GetByID)
		catMgrs.GET("/:id/details", m.handler.GetWithDetails)
		catMgrs.PUT("/:id", m.handler.Update)
		catMgrs.DELETE("/:id", m.handler.Delete)

		// Query by relationships
		catMgrs.GET("/user-category", m.handler.GetByUserAndCategory)
		catMgrs.GET("/by-user/:userID", m.handler.GetByUserID)
		catMgrs.GET("/by-category/:categoryID", m.handler.GetByCategoryID)
	}
}
