package router

import (
	"community-backend/internal/handler"
	"community-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// HistoryRoutes registers all History related routes.
func HistoryRoutes(r *gin.Engine) {
	historyAuth := r.Group("/history-posts")
	historyAuth.Use(middleware.AuthMiddleware())
	{
		historyAuth.GET("", handler.GetHistoryPostsByUserID)
	}
}