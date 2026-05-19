package repository

import (
	"github.com/seki18/lingdu-feed/internal/common"
	"github.com/seki18/lingdu-feed/internal/model"
)

// GetUserByID retrieves a single user by primary key.
func GetUserByID(id int) (model.User, error) {
	var user model.User
	err := common.DB.Get(&user, `
		SELECT id, username, password, email, following_count, follower_count, created_time
		FROM users WHERE id = $1
	`, id)
	return user, err
}

// GetUserByEmail retrieves a user by their email address.
func GetUserByEmail(email string) (model.User, error) {
	var user model.User

	err := common.DB.Get(&user, `
		SELECT id, username, email, password
		FROM users
		WHERE email = $1
	`, email)

	return user, err
}

// CreateUser inserts a new user and returns the created record.
func CreateUser(user model.User) (model.User, error) {
	err := common.DB.QueryRowx(`
		INSERT INTO users (username, password, email, created_time)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, username, email, created_time
	`, user.Username, user.Password, user.Email).
		StructScan(&user)

	return user, err
}

// UpdateUserName updates the username of an existing user and returns the updated record.
func UpdateUserName(user model.User) (model.User, error) {
	err := common.DB.QueryRowx(`
		UPDATE users
		SET username = $2
		WHERE id = $1
		RETURNING id, username, email, created_time
	`, user.ID, user.Username).
		StructScan(&user)

	return user, err
}

// UpdatePassword updates the password hash for a user by ID.
func UpdatePassword(userID int, hashedPassword string) error {
	_, err := common.DB.Exec(`
		UPDATE users SET password = $2 WHERE id = $1
	`, userID, hashedPassword)
	return err
}
