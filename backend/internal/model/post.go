package model

import "time"

// Post represents a single post row in the database.
type Post struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	Content     string    `db:"content" json:"content"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
	UpdatedTime time.Time `db:"updated_time" json:"updated_time"`
}

// CreatePostRequest is the JSON body for POST /post.
type CreatePostRequest struct {
	UserID  int    `db:"user_id" json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdatePostRequest is the JSON body for PUT /post.
type UpdatePostRequest struct {
	ID      int    `db:"id" json:"id"`
	UserID  int    `db:"user_id" json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Posts is the summary view of a post, used in post lists.
type Posts struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	Username    string    `db:"username" json:"username"`
	Title       string    `db:"title" json:"title"`
	Status	  	FeedStatus`db:"status" json:"-"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
}
