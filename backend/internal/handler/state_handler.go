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

	// First time a post becomes exposed → increment expose_count
	if req.Status >= model.StateExposed && prevStatus < model.StateExposed {
		_ = repository.IncrExposeCount(req.PostID)
	}
	// First time a post is clicked → increment view_count
	if req.Status >= model.StateClicked && prevStatus < model.StateClicked {
		_ = repository.IncrViewCount(req.PostID)
	}
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

	for i := range reqs {
		reqs[i].UserID = uid
		prevStatus := model.StateUnknown
		if existing, err := service.GetState(reqs[i]); err == nil {
			prevStatus = existing.Status
		}
		if err := service.UpsertState(reqs[i]); err != nil {
			log.Printf("[BatchUpsertState] Service error at index %d: %v", i, err)
		}
		// First time a post becomes exposed → increment expose_count
		if reqs[i].Status >= model.StateExposed && prevStatus < model.StateExposed {
			_ = repository.IncrExposeCount(reqs[i].PostID)
		}
		// First time a post is clicked → increment view_count
		if reqs[i].Status >= model.StateClicked && prevStatus < model.StateClicked {
			_ = repository.IncrViewCount(reqs[i].PostID)
		}
	}
	common.Success(c, nil)
}
