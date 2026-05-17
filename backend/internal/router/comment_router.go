package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// CommentRoutes registers all Comment related routes.
func CommentRoutes(r *gin.Engine) {
	comment := r.Group("/comments")
	{
		comment.GET("/by-post/:postId", handler.GetCommentsByPost)
		comment.GET("/:id", handler.GetCommentByID)
		comment.GET("/count/:postId", handler.GetCommentCountByPostID)
	}

	commentAuth := r.Group("/comments")
	commentAuth.Use(middleware.AuthMiddleware())
	{
		commentAuth.POST("", handler.CreateComment)
		commentAuth.DELETE("/:id", handler.DeleteCommentByID)
	}
}
