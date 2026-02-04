package http

import (
	"github.com/gin-gonic/gin"
)

type AttendanceModule struct {
	handler *AttendanceHandler
	secret  []byte
}

func NewAttendanceModule(handler *AttendanceHandler, secret []byte) *AttendanceModule {
	return &AttendanceModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *AttendanceModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	attendance := api.Group("/attendance")
	{
		attendance.POST("", m.handler.Record)
		attendance.GET("/course/:courseID", m.handler.ListByCourse)
		attendance.DELETE("/:id", m.handler.Delete)
	}
}
