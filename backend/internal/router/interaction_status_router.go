package router

import (
	"github.com/seki18/lingdu-feed/internal/handler"
	"github.com/seki18/lingdu-feed/internal/middleware"

	"github.com/gin-gonic/gin"
)

// InteractionStatusRoutes registers all InteractionStatus related routes.
func InteractionStatusRoutes(r *gin.Engine) {
	interactionStatusAuth := r.Group("/interaction-status")
	interactionStatusAuth.Use(middleware.AuthMiddleware())
	{
		interactionStatusAuth.POST("", handler.UpsetInteractionStatus)
		interactionStatusAuth.POST("/batch", handler.BatchUpsertInteractionStatus)
	}
}
