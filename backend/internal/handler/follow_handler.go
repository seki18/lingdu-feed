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

// IsFollowExist handles POST /Follows/exist (auth required). Checks if a follow exists for a given post and user.
func IsFollowExist(c *gin.Context) {
	var req model.CreateFollowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[IsFollowExist] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.FollowerID = userID.(int)
	log.Printf("[IsFollowExist] Request: following_id=%d, follower_id=%d", req.FollowingID, req.FollowerID)

	exist, err := service.IsFollowExist(req)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"exists": exist})
}

// CreateFollow handles POST /Follows (auth required). Creates a new Follow.
func CreateFollow(c *gin.Context) {
	var req model.CreateFollowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateFollow] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.FollowerID = userID.(int)
	log.Printf("[CreateFollow] Request: following_id=%d, follower_id=%d", req.FollowingID, req.FollowerID)

	// Check if already following to avoid duplicate key error
	exists, _ := service.IsFollowExist(req)
	if exists {
		common.Error(c, http.StatusConflict, common.ErrEmailExists)
		return
	}

	Follow, err := service.CreateFollow(req)
	if err != nil {
		log.Printf("[CreateFollow] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	// Sync user follow counts
	_ = repository.IncrFollowingCount(req.FollowerID)
	_ = repository.IncrFollowerCount(req.FollowingID)

	common.Success(c, Follow)
}

// DeleteFollow handles DELETE /Follows/:id (auth required). Deletes a Follow.
func DeleteFollow(c *gin.Context) {
	var req model.CreateFollowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[DeleteFollow] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.FollowerID = userID.(int)
	log.Printf("[DeleteFollow] Request: following_id=%d, follower_id=%d", req.FollowingID, req.FollowerID)

	err := service.DeleteFollow(req)
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "Follow not found") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	// Sync user follow counts
	_ = repository.DecrFollowingCount(req.FollowerID)
	_ = repository.DecrFollowerCount(req.FollowingID)

	common.Success(c, nil)
}

// GetFollowingListByFollowerID handles GET /Follows/list/following/:followerId. Returns a paginated list of users that a given user is following.
func GetFollowingListByFollowerID(c *gin.Context) {
	followerID, err := strconv.Atoi(c.Param("followerId"))
	if err != nil || followerID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	follows, total, err := service.GetFollowingListByFollowerID(followerID, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, gin.H{"follows": follows, "total": total})
}

// GetFollowerListByFollowingID handles GET /Follows/list/follower/:followingId. Returns a paginated list of followers for a given user.
func GetFollowerListByFollowingID(c *gin.Context) {
	followingID, err := strconv.Atoi(c.Param("followingId"))
	if err != nil || followingID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	follows, total, err := service.GetFollowerListByFollowingID(followingID, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, gin.H{"follows": follows, "total": total})
}
