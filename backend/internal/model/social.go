package model

import "time"

// Like represents a row in the likes table (formerly praises).
type Like struct {
	UserID      int       `db:"user_id" json:"user_id"`
	PostID      int       `db:"post_id" json:"post_id"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
}

// LikeRequest is the JSON body for like/unlike operations.
// user_id is inferred from JWT; only post_id is needed from the request.
type LikeRequest struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

// Favorite represents a row in the favorites table (formerly collections).
type Favorite struct {
	UserID      int       `db:"user_id" json:"user_id"`
	PostID      int       `db:"post_id" json:"post_id"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
}

// FavoriteRequest is the JSON body for favorite/unfavorite operations.
type FavoriteRequest struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

// CreateCommentRequest is the JSON body for creating a comment.
type CreateCommentRequest struct {
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
	ReplyID *int   `json:"reply_id"`
	Content string `json:"content"`
}

// DeleteCommentRequest is the JSON body for deleting a comment.
type DeleteCommentRequest struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}
