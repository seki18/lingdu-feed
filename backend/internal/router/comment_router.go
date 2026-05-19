package router

import (
	"github.com/seki18/lingdu-feed/internal/handler"
	"github.com/seki18/lingdu-feed/internal/middleware"

	"github.com/gin-gonic/gin"
)

// CommentRoutes registers all Comment related routes.
func CommentRoutes(r *gin.Engine) {
	commentAuth := r.Group("/comments")
	commentAuth.Use(middleware.AuthMiddleware())
	{
		commentAuth.POST("", handler.CreateComment)
		commentAuth.DELETE("/:id", handler.DeleteCommentByID)
	}

	commentPublic := r.Group("/comments")
	commentPublic.Use(middleware.SoftAuthMiddleware())
	{
		commentPublic.GET("/by-post/:post_id", handler.GetCommentsByPostID)
		commentPublic.GET("/count/:post_id", handler.GetCommentCountByPostID)
	}
}
