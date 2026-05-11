package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	user := r.Group("/users")
	{
		user.GET("/:id", handler.GetUserByID)
	}

	auth_valid := r.Group("/users")
	auth_valid.Use(middleware.AuthMiddleware())
	{
		auth_valid.GET("/me", handler.Me)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.CreateUser)
		auth.POST("/login", handler.Login)
	}
}
