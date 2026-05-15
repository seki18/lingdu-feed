package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes registers all user and auth related routes.
func UserRoutes(r *gin.Engine) {
	user := r.Group("/users")
	{
		user.GET("/:id", handler.GetUserByID)
	}

	authValid := r.Group("/users")
	authValid.Use(middleware.AuthMiddleware())
	{
		authValid.GET("/me", handler.Me)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.CreateUser)
		auth.POST("/login", handler.Login)
	}
}
