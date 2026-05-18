package model

import "time"

// Post represents a single post row in the database, including counts.
type Post struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	Username        string    `db:"username" json:"username"`
	Title           string    `db:"title" json:"title"`
	Content         string    `db:"content" json:"content"`
	PraiseCount     int       `db:"praise_count" json:"praise_count"`
	CommentCount    int       `db:"comment_count" json:"comment_count"`
	CollectionCount int       `db:"collection_count" json:"collection_count"`
	ViewCount       int       `db:"view_count" json:"view_count"`
	CreatedTime     time.Time `db:"created_time" json:"created_time"`
	UpdatedTime     time.Time `db:"updated_time" json:"updated_time"`
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

// PostDetail is the full response for a single post, including content,
// interaction status, and comments.
type PostDetail struct {
	Post         Post      `json:"post"`
	HasPraised   bool      `json:"has_praised"`
	HasCollected bool      `json:"has_collected"`
	Comments     []Comment `json:"comments"`
}
