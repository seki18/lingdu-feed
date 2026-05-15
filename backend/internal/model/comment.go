package model

import "time"

// Comment represents a single Comment row in the database.
type Comment struct {
	ID            int       `db:"id" json:"id"`
	PostID        int       `db:"post_id" json:"post_id"`
	UserID        int       `db:"user_id" json:"user_id"`
	Username      string    `db:"username" json:"username"`
	ReplyID       *int      `db:"reply_id" json:"reply_id"`
	ReplyUsername *string   `db:"reply_username" json:"reply_username"`
	Content       string    `db:"content" json:"content"`
	CreatedTime   time.Time `db:"created_time" json:"created_time"`
}

// CreateCommentRequest is the JSON body for POST /comments.
type CreateCommentRequest struct {
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
	ReplyID *int   `json:"reply_id"`
	Content string `json:"content"`
}
