package model

import "time"

// User represents a user row in the database. Password is omitted from JSON output.
type User struct {
	ID             int       `db:"id" json:"id"`
	Username       string    `db:"username" json:"username"`
	Password       string    `db:"password" json:"-"`
	Email          string    `db:"email" json:"email"`
	FollowingCount int       `db:"following_count" json:"following_count"`
	FollowerCount  int       `db:"follower_count" json:"follower_count"`
	IsFollowing    bool      `db:"is_following" json:"is_following"`
	CreatedTime    time.Time `db:"created_time" json:"created_time"`
}

// CreateUserRequest is the JSON body for POST /auth/register.
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UpdateUserRequest is the JSON body for PUT /users.
type UpdateUserRequest struct {
	ID       int    `db:"id" json:"id"`
	Username string `json:"username"`
}
