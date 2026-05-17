package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// CollectionRoutes registers all Collection related routes.
func CollectionRoutes(r *gin.Engine) {
	collections := r.Group("/Collections")
	{
		collections.GET("/count/:postId", handler.GetCollectionCountByPostID)
	}

	collectionsAuth := r.Group("/Collections")
	collectionsAuth.Use(middleware.AuthMiddleware())
	{
		collectionsAuth.GET("", handler.GetCollectionByUserID)
		collectionsAuth.POST("/exist", handler.IsCollectionExist)
		collectionsAuth.POST("", handler.CreateCollection)
		collectionsAuth.DELETE("", handler.DeleteCollection)
	}
}
