package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/service"
	"github.com/seki18/lingdu-feed/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetRecommendPosts handles GET /feed/recommend. Returns recommended posts
// with cursor-based pagination (recommend_score, recommend_id, recent_time, following_time).
func GetRecommendPosts(c *gin.Context) {
	uid := utils.GetSoftUserID(c)
	var cursors service.FeedCursors
	if rs := c.Query("recommend_score"); rs != "" {
		cursors.RecommendScore, _ = strconv.ParseFloat(rs, 64)
	}
	if ri := c.Query("recommend_id"); ri != "" {
		cursors.RecommendID, _ = strconv.Atoi(ri)
	}
	if rt := c.Query("recent_time"); rt != "" {
		cursors.RecentTime, _ = time.Parse(time.RFC3339, rt)
	}
	if ft := c.Query("following_time"); ft != "" {
		cursors.FollowingTime, _ = time.Parse(time.RFC3339, ft)
	}
	log.Printf("[GetRecommendPosts] user_id=%d cursors=%+v", uid, cursors)
	posts, newCursors, err := service.GetRecommendPosts(uid, cursors)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"posts": posts, "cursors": newCursors})
}

// GetFollowingPosts handles GET /feed/following. Returns posts from
// users that the current user follows, with cursor-based pagination.
func GetFollowingPosts(c *gin.Context) {
	uid := utils.GetAuthUserID(c)
	cursor, _ := time.Parse(time.RFC3339, c.DefaultQuery("cursor", ""))
	log.Printf("[GetFollowingPosts] user_id=%d cursor=%v", uid, cursor)
	posts, newCursor, err := service.GetFollowingPosts(uid, cursor)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"posts": posts, "cursor": newCursor.Format(time.RFC3339)})
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
