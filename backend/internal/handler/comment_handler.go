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

// CreateComment handles POST /comments (auth required). Creates a new Comment.
func CreateComment(c *gin.Context) {
	var req model.CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateComment] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[CreateComment] Request: post_id=%d user_id=%d content=%q reply_id=%v", req.PostID, req.UserID, req.Content, req.ReplyID)

	comment, err := service.CreateComment(req)
	if err != nil {
		log.Printf("[CreateComment] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	// Sync post comment_count
	_ = repository.IncrCommentCount(req.PostID)

	common.Success(c, comment)
}

// DeleteCommentByID handles DELETE /comments/:id (auth required). Deletes a comment and its replies.
func DeleteCommentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	req := model.DeleteCommentRequest{
		PostID: id,
		UserID: userID.(int),
	}
	err = service.DeleteCommentByID(req)
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "comment not found") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	// Sync post comment_count
	_ = repository.DecrCommentCount(req.PostID)

	common.Success(c, nil)
}
