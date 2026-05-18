package handler

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"community-backend/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	Follow, err := service.CreateFollow(req)
	if err != nil {
		log.Printf("[CreateFollow] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, Follow)
}

// DeleteFollow handles DELETE /Follows/:id (auth required). Deletes a Follow and its replies.
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
	common.Success(c, nil)
}

// GetFollowingCountByFollowerID handles GET /Follows/count/following/:followerId. Returns the total number of Follows for a given follower.
func GetFollowingCountByFollowerID(c *gin.Context) {
	followerID, err := strconv.Atoi(c.Param("followerId"))
	if err != nil || followerID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	count, err := service.GetFollowingCountByFollowerID(followerID)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, gin.H{"count": count})
}

// GetFollowerCountByFollowingID handles GET /Follows/count/follower/:followingId. Returns the total number of followers for a given user.
func GetFollowerCountByFollowingID(c *gin.Context) {
	followingID, err := strconv.Atoi(c.Param("followingId"))
	if err != nil || followingID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	count, err := service.GetFollowerCountByFollowingID(followingID)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, gin.H{"count": count})
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