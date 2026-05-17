package model

import "time"

// Praise represents a single Praise row in the database.
type Praise struct {
	UserID        int       `db:"user_id" json:"user_id"`
	PostID        int       `db:"post_id" json:"post_id"`
	CreatedTime   time.Time `db:"created_time" json:"created_time"`
}

// CreatePraiseRequest is the JSON body for POST /Praises.
type CreatePraiseRequest struct {
	UserID  int    `json:"user_id"`
	PostID  int    `json:"post_id"`
}
