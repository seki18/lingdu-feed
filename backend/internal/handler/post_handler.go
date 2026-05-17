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
		log.Printf("[GetPostsByUserID] Request: user_id=%d", userID)
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		posts, total, err := service.GetPostsByUserID(userID, page, pageSize)
		if err != nil {
			common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
			return
		}
		common.Success(c, gin.H{"items": posts, "total": total, "page": page, "page_size": pageSize})
		return
	}

	requestType := c.Query("request_type")
	if requestType == "" {
		requestType = "subsequent"
	}

	excludeIDs := []int{}
	if currentIDs := c.Query("current_ids"); currentIDs != "" {
		for _, idStr := range strings.Split(currentIDs, ",") {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err != nil || id <= 0 {
				common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
				return
			}
			excludeIDs = append(excludeIDs, id)
		}
	}

	log.Printf("[GetRecentPosts] Request: request_type=%s exclude_ids=%v", requestType, excludeIDs)
	posts, err := service.GetRecentPosts(requestType, excludeIDs)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, posts)
}

// DeletePostByID handles DELETE /post/:id. Deletes a post or returns 404 if not found.
func DeletePostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Printf("[DeletePostByID] Request: id=%d", id)
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
