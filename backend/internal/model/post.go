package model

import "time"

// Post represents a single post row in the database (metadata only; stats are in PostStats).
type Post struct {
	ID          int        `db:"id" json:"id"`
	UserID      int        `db:"user_id" json:"user_id"`
	Username    string     `db:"username" json:"username"`
	Title       string     `db:"title" json:"title"`
	Content     string     `db:"content" json:"content"`
	Stats       *PostStats `json:"stats,omitempty"`
	CreatedTime time.Time  `db:"created_time" json:"created_time"`
	UpdatedTime time.Time  `db:"updated_time" json:"updated_time"`
}

// CreatePostRequest is the JSON body for POST /posts.
type CreatePostRequest struct {
	UserID  int    `db:"user_id" json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdatePostRequest is the JSON body for PUT /posts/:id.
type UpdatePostRequest struct {
	ID      int    `db:"id" json:"id"`
	UserID  int    `db:"user_id" json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
