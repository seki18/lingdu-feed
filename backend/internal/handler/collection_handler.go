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

// GetCollectionByUserID handles GET /Collections (auth required). Returns paginated collections for the current user.
func GetCollectionByUserID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	collections, total, err := service.GetCollectionByUserID(userID.(int), page, pageSize)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"items": collections, "total": total, "page": page, "page_size": pageSize})
}

// IsCollectionExist handles POST /Collections/exist (auth required). Checks if a praise exists for a given post and user.
func IsCollectionExist(c *gin.Context) {
	var req model.CreateCollectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[IsCollectionExist] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[IsCollectionExist] Request: post_id=%d, user_id=%d", req.PostID, req.UserID)

	exist, err := service.IsCollectionExist(req)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"exists": exist})
}

// CreateCollection handles POST /Collections (auth required). Creates a new Collection.
func CreateCollection(c *gin.Context) {
	var req model.CreateCollectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateCollection] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[CreateCollection] Request: post_id=%d, user_id=%d", req.PostID, req.UserID)

	Collection, err := service.CreateCollection(req)
	if err != nil {
		log.Printf("[CreateCollection] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, Collection)
}

// DeleteCollection handles DELETE /Collections/:id (auth required). Deletes a Collection and its replies.
func DeleteCollection(c *gin.Context) {
	var req model.CreateCollectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[DeleteCollection] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[DeleteCollection] Request: post_id=%d, user_id=%d", req.PostID, req.UserID)

	err := service.DeleteCollection(req)
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "Collection not found") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, nil)
}

// GetCollectionCountByPostID handles GET /Collections/count/:postId. Returns the total number of Collections for a given post.
func GetCollectionCountByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	count, err := service.GetCollectionCountByPostID(postID)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, gin.H{"count": count})
}
