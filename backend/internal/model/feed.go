package model

import "time"

// FeedItem is the summary view of a post, used in all feed endpoints.
type FeedItem struct {
	ID            int        `db:"id" json:"id"`
	UserID        int        `db:"user_id" json:"user_id"`
	Username      string     `db:"username" json:"username"`
	Title         string     `db:"title" json:"title"`
	FirstImageURL *string    `db:"first_image_url" json:"first_image_url,omitempty"`
	Stats         *PostStats `json:"stats,omitempty"`
	CreatedTime   time.Time  `db:"created_time" json:"created_time"`
}
