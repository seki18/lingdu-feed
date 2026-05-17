package handler

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"community-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserByID handles GET /users/:id. Returns user details by ID.
func GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	users, err := service.GetUserByID(id)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	common.Success(c, users)
}

// CreateUser handles POST /auth/register. Creates a new user account.
func CreateUser(c *gin.Context) {
	var req model.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}

	user, err := service.CreateUser(req)
	if err != nil {
		switch err {
		case common.ErrEmailExists:
			common.Error(c, http.StatusConflict, common.ErrEmailExists)

		default:
			common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		}
		return
	}

	common.Success(c, user)
}

// Login handles POST /auth/login. Authenticates user and returns JWT token.
func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&req)

	token, err := service.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.Error(c, http.StatusNotFound, common.ErrUserNotFound)
		case common.ErrPasswordError:
			common.Error(c, http.StatusUnauthorized, common.ErrPasswordError)
		default:
			common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		}
		return
	}
	user, err := service.GetUserByEmail(req.Email)
	if err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}

	common.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Me handles GET /users/me (auth required). Returns the current user's ID.
func Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	common.Success(c, userID)
}

// UpdateUser handles PUT /users (auth required). Updates the current user's username.
func UpdateUsername(c *gin.Context) {
	var req model.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, http.StatusBadRequest, common.ErrInvalidParam.WithErr(err))
		return
	}
	userID, _ := c.Get("user_id")
	req.ID = userID.(int)
	log.Printf("[UpdateUsername] Request: user_id=%d, new_username=%s", req.ID, req.Username)

	user, err := service.UpdateUsername(req)
	if err != nil {
		log.Printf("[UpdateUsername] Service error: %v", err)
		common.Error(c, http.StatusInternalServerError, common.ErrInternalParam.WithErr(err))
		return
	}
	common.Success(c, user)
}
