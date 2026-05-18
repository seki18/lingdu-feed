package router

import (
	"github.com/seki18/lingdu-feed/internal/handler"
	"github.com/seki18/lingdu-feed/internal/middleware"

	"github.com/gin-gonic/gin"
)

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
