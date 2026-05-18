package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/service"

	"github.com/gin-gonic/gin"
)

// CreatePost handles POST /post (auth required). Creates a new post.
func CreatePost(c *gin.Context) {
	var req model.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[CreatePost] Request: user_id=%d, title=%s, content=%s", req.UserID, req.Title, req.Content)

	post, err := service.CreatePost(req)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, post)
}

// UpdatePost handles PUT /post (auth required). Updates an existing post.
func UpdatePost(c *gin.Context) {
	var req model.UpdatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	log.Printf("[UpdatePost] Request: id=%d, user_id=%d, title=%s, content=%s", req.ID, req.UserID, req.Title, req.Content)
	post, err := service.UpdatePost(req)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, post)
}

// DeletePostByID handles DELETE /post/:id. Deletes a post or returns 404 if not found.
func DeletePostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Printf("[DeletePostByID] Request: id=%d", id)
	err := service.DeletePostByID(int64(id))
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			common.Error(c, http.StatusNotFound, common.ErrPostNotFound)
			return
		}
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, nil)
}

// GetPostsByUserID handles GET /posts/:user_id. Returns all posts by a user with pagination.
func GetPostsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	log.Printf("[GetPostsByUserID] Request: user_id=%d", userID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	posts, total, err := service.GetPostsByUserID(userID, page, pageSize)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.SuccessPaginated(c, posts, total, page, pageSize)
}

// GetPostDetail handles GET /posts/:id. Returns post content, interaction status
// (has_praised, has_collected), and comments in a single response.
func GetPostDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	uid := userID.(int)

	detail, err := service.GetPostDetail(id, uid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			common.Error(c, http.StatusNotFound, common.ErrPostNotFound)
			return
		}
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}

	common.Success(c, detail)
}
