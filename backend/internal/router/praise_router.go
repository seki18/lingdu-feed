package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// PraiseRoutes registers all Praise related routes.
func PraiseRoutes(r *gin.Engine) {
	praise := r.Group("/Praises")
	{
		praise.GET("/count/:postId", handler.GetPraiseCountByPostID)
		praise.GET("/:id", handler.GetPraiseByID)
	}

	praiseAuth := r.Group("/Praises")
	praiseAuth.Use(middleware.AuthMiddleware())
	{
		praiseAuth.POST("/exist", handler.IsPraiseExist)
		praiseAuth.POST("", handler.CreatePraise)
		praiseAuth.DELETE("", handler.DeletePraise)
	}
}
