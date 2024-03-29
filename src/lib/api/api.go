package api

import (
	"lib/messages"
	"lib/tokens"
	"lib/users"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Services struct {
	Users    users.Service
	Tokens   tokens.Service
	Messages messages.Service
}

func Register(r *gin.Engine, services Services) {
	binding.Validator = NewValidator()

	r.NoRoute(NotFoundHandler())

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(ErrorMiddleware())

	v1 := r.Group("/api/v1")

	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", RegisterHandler(services.Users, services.Tokens))
		authGroup.POST("/login", LoginHandler(services.Users, services.Tokens))
		authGroup.GET("/ws", AuthMiddleware(services.Tokens), WsAuthHandler(services.Tokens))
	}

	profileGroup := v1.Group("/profile", AuthMiddleware(services.Tokens))
	{
		profileGroup.POST("/edit", ProfileEditHandler(services.Users, services.Tokens))
	}

	messagesGroup := v1.Group("/messages")
	{
		messagesGroup.GET("", AuthMiddleware(services.Tokens), MessagesListHandler(services.Messages))
		messagesGroup.POST("", AuthMiddleware(services.Tokens), MessagePostHandler(services.Messages))

		messagesGroup.GET("/ws", MessageListWebsocketHandler(services.Messages, services.Tokens))
	}
}
