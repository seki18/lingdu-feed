package router

import (
	"github.com/seki18/lingdu-feed/internal/handler"
	"github.com/seki18/lingdu-feed/internal/middleware"

	"github.com/gin-gonic/gin"
)

// PostRoutes registers all post related routes.
func PostRoutes(r *gin.Engine) {
	postAuth := r.Group("/post")
	postAuth.Use(middleware.AuthMiddleware())
	{
		postAuth.POST("", handler.CreatePost)
		postAuth.PUT("", handler.UpdatePost)
		postAuth.DELETE("/:id", handler.DeletePostByID)
	}

	posts := r.Group("/posts")
	posts.Use(middleware.SoftAuthMiddleware())
	{
		posts.GET("/:id", handler.GetPostDetail)
	}
}
