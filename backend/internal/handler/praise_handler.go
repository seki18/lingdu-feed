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

// GetPraiseByID handles GET /Praise/:id. Returns Praise details or 404 if not found.
func GetPraiseByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	Praises, err := service.GetPraiseByID(id)
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "no rows") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, Praises)
}

// IsPraiseExist handles POST /Praises/exist (auth required). Checks if a praise exists for a given post and user.
func IsPraiseExist(c *gin.Context) {
	var req model.CreatePraiseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[IsPraiseExist] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	log.Printf("[IsPraiseExist] Request: post_id=%d", req.PostID)
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)

	exist, err := service.IsPraiseExist(req)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, gin.H{"exists": exist})
}

// CreatePraise handles POST /Praises (auth required). Creates a new Praise.
func CreatePraise(c *gin.Context) {
	var req model.CreatePraiseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreatePraise] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	log.Printf("[CreatePraise] Request: post_id=%d", req.PostID)

	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	Praise, err := service.CreatePraise(req)
	if err != nil {
		log.Printf("[CreatePraise] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

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
	log.Printf("[DeletePraise] Request: post_id=%d", req.PostID)
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
	common.Success(c, nil)
}

// GetPraiseCountByPostID handles GET /Praises/count/:postId. Returns the total number of Praises for a given post.
func GetPraiseCountByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	count, err := service.GetPraiseCountByPostID(postID)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, gin.H{"count": count})
}
