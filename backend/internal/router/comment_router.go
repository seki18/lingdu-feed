package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// CommentRoutes registers all Comment related routes.
func CommentRoutes(r *gin.Engine) {
	Comment := r.Group("/comments")
	{
		Comment.GET("/by-post/:postId", handler.GetCommentsByPost)
		Comment.GET("/:id", handler.GetCommentByID)
	}

	CommentAuth := r.Group("/comments")
	CommentAuth.Use(middleware.AuthMiddleware())
	{
		CommentAuth.POST("", handler.CreateComment)
	}
}
