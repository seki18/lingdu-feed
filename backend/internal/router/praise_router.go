package router

import (
	"github.com/seki18/lingdu-feed/internal/handler"
	"github.com/seki18/lingdu-feed/internal/middleware"

	"github.com/gin-gonic/gin"
)

// PraiseRoutes registers all Praise related routes.
func PraiseRoutes(r *gin.Engine) {
	praiseAuth := r.Group("/Praises")
	praiseAuth.Use(middleware.AuthMiddleware())
	{
		praiseAuth.POST("", handler.CreatePraise)
		praiseAuth.DELETE("", handler.DeletePraise)
	}
}
