package service

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
	"community-backend/internal/repository"
	"community-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func GetUserByID(id string) (model.User, error) {
	return repository.GetUserByID(id)
}

func CreateUser(req model.CreateUserRequest) (model.User, error) {
	exist, _ := repository.GetUserByEmail(req.Email)
	if exist.ID != 0 {
		return model.User{}, common.ErrEmailExists
	}

	// 1. 加密密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return model.User{}, err
	}

	user_create := model.User{
		Username: req.Username,
		Password: string(hashedPwd),
		Email:    req.Email,
	}

	// 2. 写入数据库
	return repository.CreateUser(user_create)
}

func Login(email, password string) (string, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", err
	}

	// 生成 JWT
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
