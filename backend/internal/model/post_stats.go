package model

import "time"

// PostStats holds the high-write counter fields for a post,
// stored in the post_stats table (1:1 with posts).
type PostStats struct {
	ID            int       `db:"id" json:"id"`
	LikeCount     int       `db:"like_count" json:"like_count"`
	CommentCount  int       `db:"comment_count" json:"comment_count"`
	FavoriteCount int       `db:"favorite_count" json:"favorite_count"`
	ViewCount     int       `db:"view_count" json:"view_count"`
	ExposeCount   int       `db:"expose_count" json:"expose_count"`
	Score         float64   `db:"score" json:"score"`
	UpdatedTime   time.Time `db:"updated_time" json:"updated_time"`
}
