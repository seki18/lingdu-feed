package handler

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"community-backend/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	users, err := service.GetUserByID(id)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	common.Success(c, users)
}

func CreateUser(c *gin.Context) {
	var req model.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	user, err := service.CreateUser(req)
	if err != nil {

		switch err {
		case common.ErrEmailExists:
			common.Error(c, http.StatusConflict, common.ErrEmailExists)

		default:
			log.Printf("CreateUser error: %v", err)
			common.Error(c, http.StatusInternalServerError, common.ErrInternalParam)
		}
		return
	}

	common.Success(c, user)
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&req)

	token, err := service.Login(req.Email, req.Password)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	user, err := service.GetUserByEmail(req.Email)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	common.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	common.Success(c, userID)
}
