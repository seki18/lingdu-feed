package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(r *gin.Engine) {
	post := r.Group("/post")
	{
		post.GET("/:id", handler.GetPostByID)
	}

	postValid := r.Group("/post")
	postValid.Use(middleware.AuthMiddleware())
	{
		postValid.POST("", handler.CreatePost)
		postValid.PUT("", handler.UpdatePost)
	}

	posts := r.Group("/posts")
	{
		posts.GET("", handler.GetRecentPosts)
	}
}
