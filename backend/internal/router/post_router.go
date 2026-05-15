package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// PostRoutes registers all post related routes.
func PostRoutes(r *gin.Engine) {
	post := r.Group("/post")
	{
		post.GET("/:id", handler.GetPostByID)
	}

	postAuth := r.Group("/post")
	postAuth.Use(middleware.AuthMiddleware())
	{
		postAuth.POST("", handler.CreatePost)
		postAuth.PUT("", handler.UpdatePost)
		postAuth.DELETE("/:id", handler.DeletePostByID)
	}

	posts := r.Group("/posts")
	{
		posts.GET("", handler.GetRecentPosts)
	}
}
