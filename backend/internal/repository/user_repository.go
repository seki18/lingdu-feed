package repository

import (
	"community-backend/internal/common"
	"community-backend/internal/model"
)

func GetUserByID(id int) (model.User, error) {
	var user model.User

	err := common.DB.Get(&user, `
		SELECT id, username, password, email, created_time
		FROM users 
		WHERE id = $1
	`, id)

	return user, err
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User

	err := common.DB.Get(&user, `
		SELECT id, username, email, password
		FROM users
		WHERE email = $1
	`, email)

	return user, err
}

func CreateUser(user model.User) (model.User, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO users (username, password, email, created_time)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, username, email, created_time
	`, user.Username, user.Password, user.Email).
		StructScan(&user)

	return user, err
}
