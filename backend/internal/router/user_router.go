package router

import (
	"github.com/seki18/lingdu-feed/internal/handler"
	"github.com/seki18/lingdu-feed/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes registers all user and auth related routes.
func UserRoutes(r *gin.Engine) {
	user := r.Group("/users")
	user.Use(middleware.SoftAuthMiddleware())
	{
		user.GET("/:id", handler.GetUserByID)
	}

	authValid := r.Group("/users")
	authValid.Use(middleware.AuthMiddleware())
	{
		authValid.GET("/me", handler.Me)
		authValid.PUT("", handler.UpdateUsername)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.CreateUser)
		auth.POST("/login", handler.Login)
	}
}
