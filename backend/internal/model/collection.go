package model

import "time"

// Collection represents a single Collection row in the database.
type Collection struct {
	UserID        int       `db:"user_id" json:"user_id"`
	PostID        int       `db:"post_id" json:"post_id"`
	CreatedTime   time.Time `db:"created_time" json:"created_time"`
}

// CreateCollectionRequest is the JSON body for POST /Collections.
type CreateCollectionRequest struct {
	UserID  int    `json:"user_id"`
	PostID  int    `json:"post_id"`
}
