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
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[CreateComment] Request: post_id=%d user_id=%d content=%q reply_id=%v", req.PostID, req.UserID, req.Content, req.ReplyID)

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

// DeleteCommentByID handles DELETE /comments/:id (auth required). Deletes a comment and its replies.
func DeleteCommentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	req := model.DeleteCommentRequest{
		PostID:  id,
		UserID:  userID.(int),
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
	common.Success(c, nil)
}

// GetCommentCountByPostID handles GET /comments/count/:postId. Returns the total number of comments for a post.
func GetCommentCountByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	count, err := service.GetCommentCountByPostID(postID)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, gin.H{"count": count})
}
