package http

import (
	"github.com/gin-gonic/gin"
)

type TeacherModule struct {
	handler *TeacherHandler
	secret  []byte
}

func NewTeacherModule(handler *TeacherHandler, secret []byte) *TeacherModule {
	return &TeacherModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *TeacherModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	teachers := api.Group("/teachers")
	{
		// Get current teacher profile (for logged in teacher)
		teachers.GET("/me", m.handler.GetMyProfile)

		// CRUD operations
		teachers.POST("", m.handler.Create)
		teachers.GET("", m.handler.List)
		teachers.GET("/:id", m.handler.GetByID)
		teachers.PUT("/:id", m.handler.Update)
		teachers.DELETE("/:id", m.handler.Delete)

		// Relationships
		teachers.GET("/:id/courses", m.handler.GetCourses)
		teachers.GET("/:id/groups", m.handler.GetGroups)

		// Get teacher by user ID
		teachers.GET("/user/:userID", m.handler.GetByUserID)
	}
}
