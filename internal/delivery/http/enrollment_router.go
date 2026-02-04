package http

import (
	"github.com/gin-gonic/gin"
)

type EnrollmentModule struct {
	handler *EnrollmentHandler
	secret  []byte
}

func NewEnrollmentModule(handler *EnrollmentHandler, secret []byte) *EnrollmentModule {
	return &EnrollmentModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *EnrollmentModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	enrollments := api.Group("/enrollments")
	{
		enrollments.POST("", m.handler.Enroll)
		enrollments.GET("/course/:courseID", m.handler.ListByCourse)
		enrollments.GET("/student/:studentID", m.handler.ListByStudent)
		enrollments.DELETE("/:id", m.handler.Remove)
	}
}
