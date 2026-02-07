package http

import (
	"github.com/gin-gonic/gin"
)

type StudentModule struct {
	handler *StudentHandler
	secret  []byte
}

func NewStudentModule(handler *StudentHandler, secret []byte) *StudentModule {
	return &StudentModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *StudentModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	students := api.Group("/students")
	{
		// Get current student profile (for logged in student)
		students.GET("/me", m.handler.GetMyProfile)

		// CRUD operations
		students.POST("", m.handler.Create)
		students.GET("", m.handler.List)
		students.GET("/:id", m.handler.GetByID)
		students.GET("/:id/details", m.handler.GetWithDetails)
		students.PUT("/:id", m.handler.Update)
		students.DELETE("/:id", m.handler.Delete)

		// Relationships
		students.GET("/:id/enrollments", m.handler.GetEnrollments)
		students.GET("/:id/groups", m.handler.GetGroups)

		// Get student by user ID
		students.GET("/user/:userID", m.handler.GetByUserID)
	}
}
