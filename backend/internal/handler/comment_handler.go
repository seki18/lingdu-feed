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

// GetCommentByID handles GET /Comment/:id. Returns Comment details or 404 if not found.
func GetCommentByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	Comments, err := service.GetCommentByID(id)
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "no rows") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, Comments)
}

// CreateComment handles POST /comments (auth required). Creates a new Comment.
func CreateComment(c *gin.Context) {
	var req model.CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateComment] JSON bind error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	log.Printf("[CreateComment] Request: post_id=%d content=%q reply_id=%v", req.PostID, req.Content, req.ReplyID)

	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	comment, err := service.CreateComment(req)
	if err != nil {
		log.Printf("[CreateComment] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, comment)
}

// GetCommentsByPost handles GET /comments/by-post/:postId. Returns all comments for a post.
func GetCommentsByPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	comments, err := service.GetCommentsByPostID(postID)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, comments)
}
