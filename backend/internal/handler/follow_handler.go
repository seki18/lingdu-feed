package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
	"github.com/seki18/lingdu-feed/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateFollow handles POST /api/users/:id/follow (auth required).
// following_id comes from URL path; follower_id from JWT — no request body needed.
func CreateFollow(c *gin.Context) {
	followingID, err := strconv.Atoi(c.Param("id"))
	if err != nil || followingID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	followerID := userID.(int)
	log.Printf("[CreateFollow] Request: following_id=%d, follower_id=%d", followingID, followerID)

	req := model.CreateFollowRequest{FollowerID: followerID, FollowingID: followingID}
	exists, _ := service.IsFollowExist(req)
	if exists {
		common.Error(c, http.StatusConflict, common.ErrEmailExists)
		return
	}
	follow, err := service.CreateFollow(req)
	if err != nil {
		log.Printf("[CreateFollow] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	_ = repository.IncrFollowingCount(followerID)
	_ = repository.IncrFollowerCount(followingID)
	common.Success(c, follow)
}

// DeleteFollow handles DELETE /api/users/:id/follow (auth required).
func DeleteFollow(c *gin.Context) {
	followingID, err := strconv.Atoi(c.Param("id"))
	if err != nil || followingID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	followerID := userID.(int)
	log.Printf("[DeleteFollow] Request: following_id=%d, follower_id=%d", followingID, followerID)

	req := model.CreateFollowRequest{FollowerID: followerID, FollowingID: followingID}
	err = service.DeleteFollow(req)
	if err != nil {
		if strings.Contains(err.Error(), "Follow not found") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	_ = repository.DecrFollowingCount(followerID)
	_ = repository.DecrFollowerCount(followingID)
	common.Success(c, nil)
}

// GetFollowingListByFollowerID handles GET /api/users/:id/following.
func GetFollowingListByFollowerID(c *gin.Context) {
	followerID, err := strconv.Atoi(c.Param("id"))
	if err != nil || followerID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	follows, total, err := service.GetFollowingListByFollowerID(followerID, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"follows": follows, "total": total})
}

// GetFollowerListByFollowingID handles GET /api/users/:id/followers.
func GetFollowerListByFollowingID(c *gin.Context) {
	followingID, err := strconv.Atoi(c.Param("id"))
	if err != nil || followingID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	follows, total, err := service.GetFollowerListByFollowingID(followingID, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"follows": follows, "total": total})
}
