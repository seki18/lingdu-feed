package handler

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"community-backend/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	posts, err := service.GetPostByID(id)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	common.Success(c, posts)
}

func CreatePost(c *gin.Context) {
	var req model.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	post, err := service.CreatePost(req)
	if err != nil {
		switch err {
		default:
			log.Printf("CreatePost error: %v", err)
			common.Error(c, http.StatusInternalServerError, common.ErrInternalParam)
		}
		return
	}

	common.Success(c, post)
}

func UpdatePost(c *gin.Context) {
	var req model.UpdatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int)
	post, err := service.UpdatePost(req)
	if err != nil {
		switch err {
		default:
			log.Printf("UpdatePost error: %v", err)
			common.Error(c, http.StatusInternalServerError, common.ErrInternalParam)
		}
		return
	}

	common.Success(c, post)
}

func GetRecentPosts(c *gin.Context) {
	posts, err := service.GetRecentPosts()
	if err != nil {
		log.Printf("GetRecentPosts error: %v", err)
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	common.Success(c, posts)
}