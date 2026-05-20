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

// ── Like ──

// CreateLike handles POST /api/posts/:id/like (auth required).
// post_id comes from the URL path; user_id from JWT — no request body needed.
func CreateLike(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	uid := userID.(int)
	log.Printf("[CreateLike] Request: post_id=%d, user_id=%d", postID, uid)

	_, err = service.CreateLike(model.LikeRequest{PostID: postID, UserID: uid})
	if err != nil {
		log.Printf("[CreateLike] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	_ = repository.IncrLikeCount(postID)
	common.Success(c, nil)
}

// DeleteLike handles DELETE /api/posts/:id/like (auth required).
func DeleteLike(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	uid := userID.(int)
	log.Printf("[DeleteLike] Request: post_id=%d, user_id=%d", postID, uid)

	err = service.DeleteLike(model.LikeRequest{PostID: postID, UserID: uid})
	if err != nil {
		if strings.Contains(err.Error(), "Like not found") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	_ = repository.DecrLikeCount(postID)
	common.Success(c, nil)
}

// ── Favorite ──

// CreateFavorite handles POST /api/posts/:id/favorite (auth required).
// post_id comes from the URL path; user_id from JWT — no request body needed.
func CreateFavorite(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	uid := userID.(int)
	log.Printf("[CreateFavorite] Request: post_id=%d, user_id=%d", postID, uid)

	_, err = service.CreateFavorite(model.FavoriteRequest{PostID: postID, UserID: uid})
	if err != nil {
		log.Printf("[CreateFavorite] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	_ = repository.IncrFavoriteCount(postID)
	common.Success(c, nil)
}

// DeleteFavorite handles DELETE /api/posts/:id/favorite (auth required).
func DeleteFavorite(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	uid := userID.(int)
	log.Printf("[DeleteFavorite] Request: post_id=%d, user_id=%d", postID, uid)

	err = service.DeleteFavorite(model.FavoriteRequest{PostID: postID, UserID: uid})
	if err != nil {
		if strings.Contains(err.Error(), "favorite not found") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	_ = repository.DecrFavoriteCount(postID)
	common.Success(c, nil)
}

// ── Comment ──

// CreateComment handles POST /api/posts/:id/comments (auth required).
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
	_ = repository.IncrCommentCount(req.PostID)
	common.Success(c, comment)
}

// DeleteCommentByID handles DELETE /api/comments/:id (auth required).
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
		if strings.Contains(err.Error(), "comment not found") {
			common.Error(c, http.StatusNotFound, common.ErrNotFound)
			return
		}
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	_ = repository.DecrCommentCount(req.PostID)
	common.Success(c, nil)
}

// GetCommentsByPostID handles GET /api/posts/:id/comments (soft auth).
func GetCommentsByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	comments, err := repository.GetCommentsByPostID(postID)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, comments)
}

// GetCommentCountByPostID handles GET /api/posts/:id/comments/count (soft auth).
func GetCommentCountByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil || postID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	var count int
	err = common.DB.Get(&count, `SELECT COUNT(*) FROM comments WHERE post_id = $1`, postID)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, map[string]int{"count": count})
}
