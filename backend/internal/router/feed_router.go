package router

import (
	"github.com/seki18/lingdu-feed/internal/handler"
	"github.com/seki18/lingdu-feed/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API routes under /api.
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// ── Auth ──
	auth := api.Group("/auth")
	{
		auth.POST("/register", handler.CreateUser)
		auth.POST("/login", handler.Login)
	}

	// ── User ──
	userSoft := api.Group("/users")
	userSoft.Use(middleware.SoftAuthMiddleware())
	{
		userSoft.GET("/:id", handler.GetUserByID)
	}

	userAuth := api.Group("/users")
	userAuth.Use(middleware.AuthMiddleware())
	{
		userAuth.PUT("/me/profile", handler.UpdateUsername)
		userAuth.PUT("/me/password", handler.ChangePassword)
	}

	// ── Feed ──
	feedSoft := api.Group("/feed")
	feedSoft.Use(middleware.SoftAuthMiddleware())
	{
		feedSoft.GET("/recommend", handler.GetRecommendPosts)
		feedSoft.GET("/users/:id", handler.GetAuthorPosts)
	}

	feedAuth := api.Group("/feed")
	feedAuth.Use(middleware.AuthMiddleware())
	{
		feedAuth.GET("/following", handler.GetFollowingPosts)
		feedAuth.GET("/history", handler.GetHistoryPosts)
		feedAuth.GET("/favorites", handler.GetFavoriteFeed)
	}

	// ── Post ──
	postSoft := api.Group("/posts")
	postSoft.Use(middleware.SoftAuthMiddleware())
	{
		postSoft.GET("/:id", handler.GetPostDetail)
	}

	postAuth := api.Group("/posts")
	postAuth.Use(middleware.AuthMiddleware())
	{
		postAuth.POST("", handler.CreatePost)
		put := postAuth.Group("/:id")
		{
			put.PUT("", handler.UpdatePost)
			put.DELETE("", handler.DeletePostByID)
			put.POST("/images", handler.AddPostImages)
		}
	}

	// ── Social (like, favorite, comment) ──
	socialAuth := api.Group("/posts")
	socialAuth.Use(middleware.AuthMiddleware())
	{
		socialAuth.POST("/:id/like", handler.CreateLike)
		socialAuth.DELETE("/:id/like", handler.DeleteLike)
		socialAuth.POST("/:id/favorite", handler.CreateFavorite)
		socialAuth.DELETE("/:id/favorite", handler.DeleteFavorite)
		socialAuth.POST("/:id/comments", handler.CreateComment)
	}
	socialSoft := api.Group("/posts")
	socialSoft.Use(middleware.SoftAuthMiddleware())
	{
		socialSoft.GET("/:id/comments", handler.GetCommentsByPostID)
	}
	api.DELETE("/comments/:id", middleware.AuthMiddleware(), handler.DeleteCommentByID)

	// ── Follow ──
	followPublic := api.Group("/users")
	{
		followPublic.GET("/:id/following", handler.GetFollowingListByFollowerID)
		followPublic.GET("/:id/followers", handler.GetFollowerListByFollowingID)
	}
	followAuth := api.Group("/users")
	followAuth.Use(middleware.AuthMiddleware())
	{
		followAuth.POST("/:id/follow", handler.CreateFollow)
		followAuth.DELETE("/:id/follow", handler.DeleteFollow)
	}

	// ── Upload ──
	api.POST("/upload", middleware.AuthMiddleware(), handler.UploadImage)

	// ── State ──
	stateAuth := api.Group("/state")
	stateAuth.Use(middleware.AuthMiddleware())
	{
		stateAuth.POST("/batch", handler.BatchUpsertState)
	}
}
