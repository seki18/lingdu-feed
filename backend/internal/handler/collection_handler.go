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
	common.SuccessPaginated(c, collections, total, page, pageSize)
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

	// Sync post collection_count
	_ = repository.IncrCollectionCount(req.PostID)

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
