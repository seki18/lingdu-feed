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
// with optional request_type query parameter and current_ids for deduplication.
func GetRecommendPosts(c *gin.Context) {
	requestType := c.DefaultQuery("request_type", "subsequent")
	excludeIDs := utils.ParseExcludeIDs(c)
	uid := utils.GetSoftUserID(c)
	log.Printf("[GetRecommendPosts] Request: request_type=%s exclude_ids=%v user_id=%d", requestType, excludeIDs, uid)
	posts, err := service.GetRecommendPosts(requestType, excludeIDs, uid)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, posts)
}

// GetFollowingPosts handles GET /feed/following. Returns posts from
// users that the current user follows.
func GetFollowingPosts(c *gin.Context) {
	requestType := c.DefaultQuery("request_type", "subsequent")
	excludeIDs := utils.ParseExcludeIDs(c)
	uid := utils.GetAuthUserID(c)
	log.Printf("[GetFollowingPosts] Request: request_type=%s exclude_ids=%v user_id=%d", requestType, excludeIDs, uid)
	posts, err := service.GetFollowingPosts(requestType, excludeIDs, uid)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, posts)
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
