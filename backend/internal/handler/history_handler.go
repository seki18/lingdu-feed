package handler

import (
	"community-backend/internal/common"
	"community-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetHistoryPostsByUserID handles GET /history-posts (auth required). Returns paginated history for the current user.
func GetHistoryPostsByUserID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	posts, total, err := service.GetHistoryPostsByUserID(userID.(int), page, pageSize)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"items": posts, "total": total, "page": page, "page_size": pageSize})
}
