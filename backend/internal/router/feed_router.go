package router

import (
	"github.com/seki18/lingdu-feed/internal/handler"
	"github.com/seki18/lingdu-feed/internal/middleware"

	"github.com/gin-gonic/gin"
)

// FeedRoutes registers all feed related routes under /feed.
func FeedRoutes(r *gin.Engine) {
	// Recent posts (soft auth — works for guests too)
	recommend := r.Group("/feed/recommend")
	recommend.Use(middleware.SoftAuthMiddleware())
	{
		recommend.GET("", handler.GetRecommendPosts)
	}

	// Following posts (auth required)
	following := r.Group("/feed/following")
	following.Use(middleware.AuthMiddleware())
	{
		following.GET("", handler.GetFollowingPosts)
	}

	// Author posts (public — shows posts by a specific user)
	author := r.Group("/feed/author")
	{
		author.GET("/:user_id", handler.GetAuthorPosts)
	}

	// History (auth required)
	history := r.Group("/feed/history")
	history.Use(middleware.AuthMiddleware())
	{
		history.GET("", handler.GetHistoryPosts)
	}

	// Collections (auth required)
	collections := r.Group("/feed/collections")
	collections.Use(middleware.AuthMiddleware())
	{
		collections.GET("", handler.GetCollectionPosts)
	}
}
