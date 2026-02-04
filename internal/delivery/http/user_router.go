package http

import (
	"github.com/gin-gonic/gin"
)

type UserModule struct {
	handler *UserHandler
	secret  []byte
}

func NewUserModule(handler *UserHandler, secret []byte) *UserModule {
	return &UserModule{
		handler: handler,
		secret:  secret,
	}
}

func (m *UserModule) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	users := api.Group("/users")
	{
		users.POST("/register", m.handler.Register)
		users.POST("/login", m.handler.Login)
		users.GET("", m.handler.List)
		users.GET("/:id", m.handler.GetByID)
		users.PATCH("/:id", m.handler.Update)
		users.DELETE("/:id", m.handler.Delete)
	}
}
