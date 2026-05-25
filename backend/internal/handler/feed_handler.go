package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/service"
	"github.com/seki18/lingdu-feed/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetRecommendPosts handles GET /feed/recommend. Returns recommended posts
// with cursor-based pagination (single cursor ID).
func GetRecommendPosts(c *gin.Context) {
	uid := utils.GetSoftUserID(c)
	requestType := c.DefaultQuery("request_type", "subsequent")
	cursorID, _ := strconv.Atoi(c.DefaultQuery("cursor", "0"))
	log.Printf("[GetRecommendPosts] requestType=%s user_id=%d cursor=%d", requestType, uid, cursorID)
	posts, newCursor, err := service.GetRecommendPosts(uid, requestType, cursorID)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"posts": posts, "cursor": newCursor})
}

// GetFollowingPosts handles GET /feed/following. Returns posts from
// users that the current user follows, with cursor-based pagination.
func GetFollowingPosts(c *gin.Context) {
	uid := utils.GetAuthUserID(c)
	cursorID, _ := strconv.Atoi(c.DefaultQuery("cursor", "0"))
	log.Printf("[GetFollowingPosts] user_id=%d cursor=%d", uid, cursorID)
	posts, newCursor, err := service.GetFollowingPosts(uid, cursorID)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"posts": posts, "cursor": newCursor})
}

// GetHistoryPosts handles GET /feed/history.
func GetHistoryPosts(c *gin.Context) {
	uid := utils.GetAuthUserID(c)
	page, pageSize := utils.ParsePagination(c)
	log.Printf("[GetHistoryPosts] user_id=%d page=%d page_size=%d", uid, page, pageSize)
	posts, total, err := service.GetHistoryPosts(uid, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.SuccessPaginated(c, posts, total, page, pageSize)
}

// GetFavoriteFeed handles GET /feed/favorites.
func GetFavoriteFeed(c *gin.Context) {
	uid := utils.GetAuthUserID(c)
	page, pageSize := utils.ParsePagination(c)
	log.Printf("[GetFavoriteFeed] user_id=%d page=%d page_size=%d", uid, page, pageSize)
	posts, total, err := service.GetFavoriteFeed(uid, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.SuccessPaginated(c, posts, total, page, pageSize)
}

// GetAuthorPosts handles GET /feed/author/:user_id. Returns author profile and authored posts in one response.
func GetAuthorPosts(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	page, pageSize := utils.ParsePagination(c)
	currentUserID := utils.GetSoftUserID(c)
	log.Printf("[GetAuthorPosts] user_id=%d current_user_id=%d page=%d page_size=%d", userID, currentUserID, page, pageSize)
	user, posts, total, err := service.GetAuthorPage(userID, currentUserID, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}

	type authorPageResponse struct {
		User     model.User       `json:"user"`
		Posts    []model.FeedItem `json:"posts"`
		Total    int              `json:"total"`
		Page     int              `json:"page"`
		PageSize int              `json:"page_size"`
	}

	common.Success(c, authorPageResponse{
		User:     user,
		Posts:    posts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
