package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// FollowRoutes registers all Follow related routes.
func FollowRoutes(r *gin.Engine) {
	follow := r.Group("/Follows")
	{
		follow.GET("/count/following/:followerId", handler.GetFollowingCountByFollowerID)
		follow.GET("/count/follower/:followingId", handler.GetFollowerCountByFollowingID)
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
