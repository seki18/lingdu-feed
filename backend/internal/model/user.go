package model

import "time"

type User struct {
	ID          int       `db:"id" json:"id"`
	Username    string    `db:"username" json:"username"`
	Password    string    `db:"password" json:"-"`
	Email       string    `db:"email" json:"email"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
