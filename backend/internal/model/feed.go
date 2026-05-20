package model

import "time"

// FeedItem is the summary view of a post, used in all feed endpoints.
type FeedItem struct {
	ID            int       `db:"id" json:"id"`
	UserID        int       `db:"user_id" json:"user_id"`
	Username      string    `db:"username" json:"username"`
	Title         string    `db:"title" json:"title"`
	LikeCount     int       `db:"like_count" json:"like_count"`
	CommentCount  int       `db:"comment_count" json:"comment_count"`
	FavoriteCount int       `db:"favorite_count" json:"favorite_count"`
	ViewCount     int       `db:"view_count" json:"view_count"`
	CreatedTime   time.Time `db:"created_time" json:"created_time"`
}
