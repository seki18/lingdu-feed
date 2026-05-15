package handler

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"community-backend/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetPostByID handles GET /post/:id. Returns post details or 404 if not found.
func GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	posts, err := service.GetPostByID(id)
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "no rows") {
			common.Error(c, http.StatusNotFound, common.ErrPostNotFound)
			return
		}
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, posts)
}

// CreatePost handles POST /post (auth required). Creates a new post.
func CreatePost(c *gin.Context) {
	var req model.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
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
	post, err := service.UpdatePost(req)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}

	common.Success(c, post)
}

// GetRecentPosts handles GET /posts. Returns the most recent posts,
// or filtered by user_id query parameter.
func GetRecentPosts(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID <= 0 {
			common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
			return
		}
		posts, err := service.GetPostsByUserID(userID)
		if err != nil {
			common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
			return
		}
		common.Success(c, posts)
		return
	}

	posts, err := service.GetRecentPosts()
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, posts)
}

// DeletePostByID handles DELETE /post/:id. Deletes a post or returns 404 if not found.
func DeletePostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := service.DeletePostByID(int64(id))
	if err != nil {
		// Distinguish "not found" from other errors
		if strings.Contains(err.Error(), "no rows") {
			common.Error(c, http.StatusNotFound, common.ErrPostNotFound)
			return
		}
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, nil)
}
