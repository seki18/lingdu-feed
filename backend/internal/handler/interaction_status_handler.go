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

// GetInteractionStatus handles GET /interaction-status/:id. Returns InteractionStatus details or 404 if not found.
func GetInteractionStatus(c *gin.Context) {
	var req model.CreateInteractionStatusRequest
	req.PostID, _ = strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[GetInteractionStatus] Request: post_id=%d user_id=%d status=%d", req.PostID, req.UserID, req.Status)

	interactionStatus, err := service.GetInteractionStatus(req)
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "no rows") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, interactionStatus)
}

// GetInteractionStatusByUserID handles GET /interaction-status. Returns InteractionStatus details or 404 if not found.
func GetInteractionStatusByUserID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	interactionStatus, err := service.GetInteractionStatusByUserID(userID.(int))
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "no rows") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, interactionStatus)
}

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
	// log.Printf("[UpsetInteractionStatus] Request: post_id=%d user_id=%d status=%d", req.PostID, req.UserID, req.Status)

	err := service.UpsertInteractionStatus(req)
	if err != nil {
		log.Printf("[UpsetInteractionStatus] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, nil)
}