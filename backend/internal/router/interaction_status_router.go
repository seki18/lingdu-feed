package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// InteractionStatusRoutes registers all InteractionStatus related routes.
func InteractionStatusRoutes(r *gin.Engine) {
	// interactionStatus := r.Group("/interaction-status")
	// {
	// 	interactionStatus.GET("/:id", handler.GetInteractionStatus)
	// 	interactionStatus.GET("", handler.GetInteractionStatusByUserID)
	// }

	interactionStatusAuth := r.Group("/interaction-status")
	interactionStatusAuth.Use(middleware.AuthMiddleware())
	{
		interactionStatusAuth.POST("", handler.UpsetInteractionStatus)
	}
}