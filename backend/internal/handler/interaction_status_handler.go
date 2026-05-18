package handler

import (
	"log"
	"net/http"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
	"github.com/seki18/lingdu-feed/internal/service"

	"github.com/gin-gonic/gin"
)

// UpsetInteractionStatus handles POST /interaction-status (auth required). Creates or updates an InteractionStatus.
func UpsetInteractionStatus(c *gin.Context) {
	var req model.CreateInteractionStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpsetInteractionStatus] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[UpsetInteractionStatus] Request: post_id=%d user_id=%d status=%d", req.PostID, req.UserID, req.Status)

	err := service.UpsertInteractionStatus(req)
	if err != nil {
		log.Printf("[UpsetInteractionStatus] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	// Only increment view count when status first reaches click level
	prevStatus := model.FeedUnknown
	if existing, err := service.GetInteractionStatus(req); err == nil {
		prevStatus = existing.Status
	}
	if req.Status >= model.FeedClick && prevStatus < model.FeedClick {
		_ = repository.IncrViewCount(req.PostID)
	}

	common.Success(c, nil)
}

// BatchUpsertInteractionStatus handles POST /interaction-status/batch (auth required).
// Accepts an array of {post_id, status} and bulk upserts them all.
func BatchUpsertInteractionStatus(c *gin.Context) {
	var reqs []model.CreateInteractionStatusRequest

	if err := c.ShouldBindJSON(&reqs); err != nil {
		log.Printf("[BatchUpsertInteractionStatus] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	uid := userID.(int)

	for i := range reqs {
		reqs[i].UserID = uid
		// Check previous status to decide whether view count should increment
		prevStatus := model.FeedUnknown
		if existing, err := service.GetInteractionStatus(reqs[i]); err == nil {
			prevStatus = existing.Status
		}
		if err := service.UpsertInteractionStatus(reqs[i]); err != nil {
			log.Printf("[BatchUpsertInteractionStatus] Service error at index %d: %v", i, err)
		}
		// Only increment view count when status first reaches click level
		if reqs[i].Status >= model.FeedClick && prevStatus < model.FeedClick {
			_ = repository.IncrViewCount(reqs[i].PostID)
		}
	}

	common.Success(c, nil)
}
