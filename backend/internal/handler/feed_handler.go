package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/seki18/lingdu-feed/internal/common"
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

// GetCollectionPosts handles GET /feed/collections.
func GetCollectionPosts(c *gin.Context) {
	uid := utils.GetAuthUserID(c)
	page, pageSize := utils.ParsePagination(c)
	log.Printf("[GetCollectionPosts] user_id=%d page=%d page_size=%d", uid, page, pageSize)
	posts, total, err := service.GetCollectionPosts(uid, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.SuccessPaginated(c, posts, total, page, pageSize)
}

// GetAuthorPosts handles GET /feed/author/:user_id. Returns posts by a specific user.
func GetAuthorPosts(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	page, pageSize := utils.ParsePagination(c)
	log.Printf("[GetAuthorPosts] user_id=%d page=%d page_size=%d", userID, page, pageSize)
	posts, total, err := service.GetAuthorPosts(userID, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.SuccessPaginated(c, posts, total, page, pageSize)
}
