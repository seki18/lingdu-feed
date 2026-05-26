package handler

import (
	"log"
	"net/http"

	"github.com/seki18/lingdu-feed/internal/cache"
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/service"

	"github.com/gin-gonic/gin"
)

// UpsertState handles POST /api/state (auth required).
func UpsertState(c *gin.Context) {
	var req model.StateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpsertState] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[UpsertState] Request: post_id=%d user_id=%d status=%d", req.PostID, req.UserID, req.Status)

	// Fetch previous state before upserting (so we can detect transitions)
	prevStatus := model.StateUnknown
	if existing, err := service.GetState(req); err == nil {
		prevStatus = existing.Status
	}

	err := service.UpsertState(req)
	if err != nil {
		log.Printf("[UpsertState] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	// First time a post becomes exposed → increment expose_count (via service → cache)
	if req.Status >= model.StateExposed && prevStatus < model.StateExposed {
		_ = service.IncrExposeCount(req.PostID)
	}
	// First time a post is clicked → increment view_count (via service → cache)
	if req.Status >= model.StateClicked && prevStatus < model.StateClicked {
		_ = service.IncrViewCount(req.PostID)
	}

	// Mark consumed in Redis SET cache (best-effort)
	_ = cache.MarkConsumed(req.UserID, []int{req.PostID})

	common.Success(c, nil)
}

// BatchUpsertState handles POST /api/state/batch (auth required).
func BatchUpsertState(c *gin.Context) {
	var reqs []model.StateRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		log.Printf("[BatchUpsertState] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	uid := userID.(int)

	consumedPostIDs := make([]int, 0, len(reqs))
	for i := range reqs {
		reqs[i].UserID = uid
		prevStatus := model.StateUnknown
		if existing, err := service.GetState(reqs[i]); err == nil {
			prevStatus = existing.Status
		}
		if err := service.UpsertState(reqs[i]); err != nil {
			log.Printf("[BatchUpsertState] Service error at index %d: %v", i, err)
		}
		// First time a post becomes exposed → increment expose_count (via service → cache)
		if reqs[i].Status >= model.StateExposed && prevStatus < model.StateExposed {
			_ = service.IncrExposeCount(reqs[i].PostID)
		}
		// First time a post is clicked → increment view_count (via service → cache)
		if reqs[i].Status >= model.StateClicked && prevStatus < model.StateClicked {
			_ = service.IncrViewCount(reqs[i].PostID)
		}
		// Track consumed posts for Redis SET cache
		if reqs[i].Status > model.StateDelivered {
			consumedPostIDs = append(consumedPostIDs, reqs[i].PostID)
		}
	}

	// Mark consumed in Redis SET cache (best-effort, batch)
	_ = cache.MarkConsumed(uid, consumedPostIDs)

	common.Success(c, nil)
}
