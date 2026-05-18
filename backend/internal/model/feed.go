package model

import "time"

// Posts is the summary view of a post, used in all feed endpoints.
type Posts struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	Username        string    `db:"username" json:"username"`
	Title           string    `db:"title" json:"title"`
	PraiseCount     int       `db:"praise_count" json:"praise_count"`
	CommentCount    int       `db:"comment_count" json:"comment_count"`
	CollectionCount int       `db:"collection_count" json:"collection_count"`
	ViewCount       int       `db:"view_count" json:"view_count"`
	CreatedTime     time.Time `db:"created_time" json:"created_time"`
}
