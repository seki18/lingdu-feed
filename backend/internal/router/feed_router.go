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
		authValid.PUT("/password", handler.ChangePassword)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.CreateUser)
		auth.POST("/login", handler.Login)
	}
}


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
		posts.POST("/batch-stats", handler.BatchGetPostStats)
	}
}

// PraiseRoutes registers all Praise related routes.
func PraiseRoutes(r *gin.Engine) {
	praiseAuth := r.Group("/Praises")
	praiseAuth.Use(middleware.AuthMiddleware())
	{
		praiseAuth.POST("", handler.CreatePraise)
		praiseAuth.DELETE("", handler.DeletePraise)
	}
}

// CollectionRoutes registers all Collection related routes.
func CollectionRoutes(r *gin.Engine) {
	collectionsAuth := r.Group("/Collections")
	collectionsAuth.Use(middleware.AuthMiddleware())
	{
		collectionsAuth.GET("", handler.GetCollectionByUserID)
		collectionsAuth.POST("", handler.CreateCollection)
		collectionsAuth.DELETE("", handler.DeleteCollection)
	}
}

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

// FollowRoutes registers all Follow related routes.
func FollowRoutes(r *gin.Engine) {
	follow := r.Group("/Follows")
	{
		follow.GET("/list/following/:followerId", handler.GetFollowingListByFollowerID)
		follow.GET("/list/follower/:followingId", handler.GetFollowerListByFollowingID)
	}

	followAuth := r.Group("/Follows")
	followAuth.Use(middleware.AuthMiddleware())
	{
		followAuth.POST("/exist", handler.IsFollowExist)
		followAuth.POST("", handler.CreateFollow)
		followAuth.DELETE("", handler.DeleteFollow)
	}
}

// InteractionStatusRoutes registers all InteractionStatus related routes.
func InteractionStatusRoutes(r *gin.Engine) {
	interactionStatusAuth := r.Group("/interaction-status")
	interactionStatusAuth.Use(middleware.AuthMiddleware())
	{
		interactionStatusAuth.POST("", handler.UpsetInteractionStatus)
		interactionStatusAuth.POST("/batch", handler.BatchUpsertInteractionStatus)
	}
}