package service

import (
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
	"github.com/seki18/lingdu-feed/internal/repository"
	"github.com/seki18/lingdu-feed/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// GetUserByID retrieves a user by their primary key ID.
func GetUserByID(id int) (model.User, error) {
	return repository.GetUserByID(id)
}

// GetUserByEmail retrieves a user by their email address.
func GetUserByEmail(email string) (model.User, error) {
	return repository.GetUserByEmail(email)
}

// CreateUser registers a new user with bcrypt-hashed password.
func CreateUser(req model.CreateUserRequest) (model.User, error) {
	exist, _ := repository.GetUserByEmail(req.Email)
	if exist.ID != 0 {
		return model.User{}, common.ErrEmailExists
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return model.User{}, err
	}

	userCreate := model.User{
		Username: req.Username,
		Password: string(hashedPwd),
		Email:    req.Email,
	}

	return repository.CreateUser(userCreate)
}

// Login verifies email/password and returns a signed JWT token.
func Login(email, password string) (string, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return "", common.ErrUserNotFound
	}

	if user.ID == 0 {
		return "", common.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", common.ErrPasswordError
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// UpdateUsername updates the username of an existing user.
func UpdateUsername(req model.UpdateUserRequest) (model.User, error) {
	userUpdate := model.User{
		ID:       req.ID,
		Username: req.Username,
	}
	return repository.UpdateUserName(userUpdate)
}
