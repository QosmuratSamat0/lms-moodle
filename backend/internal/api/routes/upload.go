// Package routes contains route registration helpers organized by domain
package routes

import (
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/upload"
	"github.com/gin-gonic/gin"
)

// UploadRouter handles upload-related routes
type UploadRouter struct {
	uploadHandler *upload.Handler
}

// NewUploadRouter creates a new upload router
func NewUploadRouter(uploadHandler *upload.Handler) *UploadRouter {
	return &UploadRouter{
		uploadHandler: uploadHandler,
	}
}

// SetupRoutes configures upload routes
func (ur *UploadRouter) SetupRoutes(api *gin.RouterGroup) {
	uploads := api.Group("/uploads")
	{
		// Public info (authenticated)
		uploads.GET("/allowed-types", ur.uploadHandler.GetAllowedTypes)

		// File operations
		uploads.POST("", ur.uploadHandler.Upload)
		uploads.POST("/url", ur.uploadHandler.UploadFromURL)
		uploads.GET("/my", ur.uploadHandler.ListMyUploads)
		uploads.GET("/:id", ur.uploadHandler.GetByID)
		uploads.DELETE("/:id", ur.uploadHandler.Delete)
		uploads.GET("/reference/:type/:reference_id", ur.uploadHandler.ListByReference)
	}
}
