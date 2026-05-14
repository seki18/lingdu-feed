package model

import "time"

type Post struct {
	ID          int       `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"user_id"`
	Title    string    `db:"title" json:"title"`
	Content       string    `db:"content" json:"content"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
	UpdatedTime time.Time `db:"updated_time" json:"updated_time"`
}

type CreatePostRequest struct {
	UserID    int    `db:"user_id" json:"user_id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	ID          int       `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"user_id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

type Posts struct {
	ID          int       `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"user_id"`
	Username    string    `db:"username" json:"username"`
	Title    string    `db:"title" json:"title"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
}
