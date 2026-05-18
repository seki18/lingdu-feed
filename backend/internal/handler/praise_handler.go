package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
	"github.com/seki18/lingdu-feed/internal/service"

	"github.com/gin-gonic/gin"
)

// CreatePraise handles POST /Praises (auth required). Creates a new Praise.
func CreatePraise(c *gin.Context) {
	var req model.CreatePraiseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreatePraise] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[CreatePraise] Request: post_id=%d, user_id=%d", req.PostID, req.UserID)

	Praise, err := service.CreatePraise(req)
	if err != nil {
		log.Printf("[CreatePraise] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	// Sync post praise_count
	_ = repository.IncrPraiseCount(req.PostID)

	common.Success(c, Praise)
}

// DeletePraise handles DELETE /Praises/:id (auth required). Deletes a Praise and its replies.
func DeletePraise(c *gin.Context) {
	var req model.CreatePraiseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[DeletePraise] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[DeletePraise] Request: post_id=%d, user_id=%d", req.PostID, req.UserID)

	err := service.DeletePraise(req)
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "Praise not found") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	// Sync post praise_count
	_ = repository.DecrPraiseCount(req.PostID)

	common.Success(c, nil)
}
